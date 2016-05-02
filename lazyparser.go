package main

import (
	"strings"
)

type metadata struct {
	author   string
	title    string
	diffName string
	creator  string
}

// lazyParser gets basic beatmap info from a beatmap file. Possibly using the
// laziest methods you've ever seen.
func lazyParser(dataB []byte) metadata {
	var md metadata
	var sec string
	data := string(dataB)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\r")
		parts := strings.Split(line, ":")
		var breakAll bool
		switch {
		case len(parts) == 1 && len(parts[0]) > 0:
			if parts[0][0] == '[' && parts[0][len(parts[0])-1] == ']' {
				// if we're switching to a new section and we were in Metadata
				if sec == "[Metadata]" {
					breakAll = true
				}
				sec = parts[0]
			}
		case len(parts) == 2 && sec == "[Metadata]":
			switch parts[0] {
			case "Title":
				md.title = strings.Trim(parts[1], " ")
			case "Artist":
				md.author = strings.Trim(parts[1], " ")
			case "Creator":
				md.creator = strings.Trim(parts[1], " ")
			case "Version":
				md.diffName = strings.Trim(parts[1], " ")
			}
		}
		if breakAll {
			break
		}
	}
	return md
}
