package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

// BeatmapPOST submits a custom (unsubmitted) beatmap to OaaS,
// so that replays with that beatmap can be calculated.
func BeatmapPOST(c *gin.Context) {
	file, _, err := c.Request.FormFile("beatmap")
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}
	rawData, err := ioutil.ReadAll(file)
	defer file.Close()
	// max 512 kb
	if len(rawData) > (1024 * 512) {
		c.JSON(413, baseResponse{false, "Max .osu file size 512kb"})
		return
	}

	if string(rawData[:15]) != "osu file format" {
		c.JSON(400, baseResponse{false, "Not a valid .osu file (must begin with 'osu file format')"})
		return
	}

	md5sum := fmt.Sprintf("%x", md5.Sum(rawData))

	if beatmapIDIfExists(md5sum) != 0 {
		c.JSON(400, baseResponse{false, "beatmap already exists in db"})
		return
	}

	f, err := os.Create("maps/" + md5sum + ".osu")
	if err != nil {
		c.Error(err)
		c.JSON(500, err500)
		return
	}
	f.Write(rawData)

	beatmapInfo := lazyParser(rawData)

	db.Exec("INSERT INTO beatmaps(md5, author, title, diff_name, creator) VALUES (?, ?, ?, ?, ?)",
		md5sum, beatmapInfo.author, beatmapInfo.title, beatmapInfo.diffName, beatmapInfo.creator)

	c.JSON(200, baseResponse{Ok: true})
	return
}
