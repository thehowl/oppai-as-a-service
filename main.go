package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thehowl/conf"
	"gopkg.in/thehowl/go-osuapi.v1"
)

var db *sql.DB
var api *osuapi.Client
var cf confSt

func main() {
	err := conf.Load(&cf, "oaas.conf")
	if err == conf.ErrNoFile {
		cf.Workers = 8
		err = conf.Export(&cf, "oaas.conf")
		if err != nil {
			panic(err)
		}
		fmt.Println("generated sample oaas.conf, please set the values appropriately")
		return
	}

	// start db connection
	db, err = sql.Open("mysql", cf.DSN)
	if err != nil {
		fmt.Println(err)
		return
	}

	// start api client
	api = osuapi.NewClient(cf.APIKey)
	err = api.Test()
	if err != nil {
		fmt.Println(err)
		return
	}

	// start a few workers
	if cf.Workers == 0 {
		cf.Workers = 8
	}
	for i := 0; i < cf.Workers; i++ {
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

type confSt struct {
	DSN     string `description:"go-mysql-driver dsn"`
	APIKey  string `description:"osu api key"`
	Workers int
}
