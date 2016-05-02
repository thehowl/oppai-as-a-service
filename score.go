package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
	"github.com/thehowl/osrtool/osr"
	"gopkg.in/thehowl/go-osuapi.v1"
)

type scorePostResponse struct {
	baseResponse
	ScoreID int `json:"score_id"`
}

// ScorePOST creates a new score in the DB, from its osr file.
func ScorePOST(c *gin.Context) {
	file, _, err := c.Request.FormFile("replay")
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}
	rawData, err := ioutil.ReadAll(file)
	defer file.Close()
	// max 1 mb
	if len(rawData) > (1024 * 1024) {
		c.JSON(413, baseResponse{false, "Max replay size 1mb"})
		return
	}

	// check whether score is already in db
	md5sum := fmt.Sprintf("%x", md5.Sum(rawData))
	scoreID := getScoreIDIfExists(md5sum)
	if scoreID != 0 {
		s := scorePostResponse{}
		s.Ok = true
		s.ScoreID = scoreID
		c.JSON(200, s)
		return
	}

	// get osr file data
	rep, err := osr.Unmarshal(bytes.NewBuffer(rawData))
	if err != nil {
		c.Error(err)
		c.JSON(500, baseResponse{false, "An error occurred while trying to parse your replay. Are you sure it's an .osr file?"})
		return
	}

	// kick out all non-std replays
	if rep.GameMode != osr.ModeStandard {
		c.JSON(400, baseResponse{false, "Only osu!standard PP calculation is supported."})
		return
	}

	// remove akerino attempts from beatmap hashes, and generate file name for beatmap
	rep.BeatmapHash = sanitize.BaseName(rep.BeatmapHash)
	destFileName := fmt.Sprintf("maps/%s.osu", rep.BeatmapHash)

	// check whether we already have the beatmap
	bid := beatmapIDIfExists(rep.BeatmapHash)
	if bid == 0 {
		// we don't have the beatmap. retrieve the data from the osu! API
		beatmaps, err := api.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapHash: rep.BeatmapHash,
		})
		// make sure a beatmap is returned
		if err != nil || len(beatmaps) == 0 {
			if err != nil {
				c.Error(err)
			}
			c.JSON(404, baseResponse{false, "Looked for the beatmap on the osu! website, but it could not be found. If you have the .osu file, please upload it."})
			return
		}
		// get .osu file
		websiteResp, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/osu/%d", beatmaps[0].BeatmapID))
		if err != nil {
			c.Error(err)
			c.JSON(500, err500)
			return
		}
		// create destination file
		destFile, err := os.Create(destFileName)
		if err != nil {
			c.Error(err)
			c.JSON(500, err500)
			return
		}
		defer destFile.Close()
		// copy all data from .osu file retrieved with the osu! website to the local file
		_, err = io.Copy(destFile, websiteResp.Body)
		if err != nil {
			c.Error(err)
			c.JSON(500, err500)
			return
		}
		// create beatmap in db
		insertRes, err := db.Exec("INSERT INTO beatmaps(md5, author, title, diff_name, creator) VALUES (?, ?, ?, ?, ?)",
			rep.BeatmapHash, beatmaps[0].Artist, beatmaps[0].Title, beatmaps[0].DiffName, beatmaps[0].Creator)
		if err != nil {
			c.Error(err)
			c.JSON(500, err500)
			return
		}
		// get beatmap id in db
		lid, err := insertRes.LastInsertId()
		if err != nil {
			c.Error(err)
			c.JSON(500, err500)
			return
		}
		// return beatmap id
		bid = int(lid)
	}
	// calculate score accuracy
	acc := calculateAccuracy(int(rep.Hit300), int(rep.Hit100), int(rep.Hit50), int(rep.HitGeki), int(rep.HitKatu), int(rep.HitMiss), 0)
	// insert score into table
	lidRaw, err := db.Exec("INSERT INTO scores (replay_md5, beatmap_id, player, accuracy, mods, max_combo, misses) VALUES (?, ?, ?, ?, ?, ?, ?)",
		string(md5sum[:]), bid, rep.Player, acc, int(rep.Mods), int(rep.MaxCombo), rep.HitMiss)
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}
	// get score id
	lid, err := lidRaw.LastInsertId()
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}

	// enqueue pp calculation task
	tasks <- oppaiTask{
		ScoreID:  int(lid),
		FilePath: destFileName,
		Accuracy: acc,
		Mods:     osuapi.Mods(rep.Mods),
		MaxCombo: int(rep.MaxCombo),
		Misses:   int(rep.HitMiss),
	}

	// respond
	s := scorePostResponse{}
	s.Ok = true
	s.ScoreID = int(lid)
	c.JSON(200, s)
}

func getScoreIDIfExists(md5 string) int {
	// suppress errors because yolo
	var scoreID int
	db.QueryRow("SELECT id FROM scores WHERE replay_md5 = ? LIMIT 1", string(md5)).Scan(&scoreID)
	return scoreID
}

func beatmapIDIfExists(beatmapMD5 string) int {
	var id int
	db.QueryRow("SELECT id FROM beatmaps WHERE md5 = ? LIMIT 1", beatmapMD5).Scan(&id)
	return id
}

type scoreData struct {
	baseResponse
	Calculated int `json:"calculated"`
	Score      struct {
		Player   string  `json:"player"`
		Accuracy float64 `json:"accuracy"`
		Mods     int     `json:"mods"`
		ModsStr  string  `json:"mods_str"`
		PP       float64 `json:"pp"`
	} `json:"score"`
	Beatmap struct {
		Author   string `json:"author"`
		Title    string `json:"title"`
		DiffName string `json:"diff_name"`
		Creator  string `json:"creator"`
	} `json:"beatmap"`
}

// ScoreGET retrives data of a score knowing its ID.
func ScoreGET(c *gin.Context) {
	// prepare scoreID
	if c.Query("id") == "" {
		c.JSON(400, baseResponse{false, "Please provide a score ID (param id)"})
		return
	}
	scoreID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Error(err)
		c.JSON(400, baseResponse{false, "Please provide a valid number as the score ID"})
		return
	}

	// get a shitload of data
	var sd scoreData
	err = db.QueryRow(`
	SELECT 
		scores.player, scores.accuracy, scores.mods,
		scores.calculated, scores.total_pp,
		beatmaps.author, beatmaps.title, beatmaps.diff_name, beatmaps.creator
	FROM scores
	LEFT JOIN beatmaps
		ON scores.beatmap_id = beatmaps.id
	WHERE scores.id = ?
	LIMIT 1`, scoreID).Scan(
		&sd.Score.Player, &sd.Score.Accuracy, &sd.Score.Mods,
		&sd.Calculated, &sd.Score.PP,
		&sd.Beatmap.Author, &sd.Beatmap.Title, &sd.Beatmap.DiffName, &sd.Beatmap.Creator,
	)
	if err == sql.ErrNoRows {
		c.JSON(404, baseResponse{false, "That score could not be found!"})
	}
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}
	sd.Ok = true
	c.JSON(200, sd)
}

// ScoreSubmitGET allows for score submission of a score without having to give a replay.
func ScoreSubmitGET(c *gin.Context) {

}
