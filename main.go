package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/thehowl/go-osuapi.v1"
)

var db *sql.DB
var api *osuapi.Client

func main() {
	dsn := os.Getenv("DSN")
	apiKey := os.Getenv("OSU_API_KEY")
	if dsn == "" {
		fmt.Println("Please set a DSN for connecting to the database")
		return
	}
	if apiKey == "" {
		fmt.Println("Please set an OSU_API_KEY")
		return
	}

	// start db connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// start api client
	api = osuapi.NewClient(apiKey)
	err = api.Test()
	if err != nil {
		fmt.Println(err)
		return
	}

	// start a few workers
	for i := 0; i < 8; i++ {
		go Worker()
	}

	// start webserver
	app := gin.Default()

	app.LoadHTMLFiles("static/page.html")

	app.POST("/api/v1/beatmap", BeatmapPOST)
	app.POST("/api/v1/score", ScorePOST)
	app.GET("/api/v1/score", ScoreGET)
	app.GET("/api/v1/score/submit", ScoreSubmitGET)

	app.Any("/", Frontend)

	app.Static("/static", "static")

	app.Run(":42043")
}
