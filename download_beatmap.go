package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadOsuFile(beatmapID int, destFileName string) error {
	websiteResp, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/osu/%d", beatmapID))
	if err != nil {
		return err
	}
	// create destination file
	destFile, err := os.Create(destFileName)
	if err != nil {
		return err
	}
	defer destFile.Close()
	// copy all data from .osu file retrieved with the osu! website to the local file
	_, err = io.Copy(destFile, websiteResp.Body)
	if err != nil {
		return err
	}
	return nil
}
