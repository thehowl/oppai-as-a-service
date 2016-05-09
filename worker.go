package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/thehowl/go-osuapi.v1"
)

type oppaiTask struct {
	ScoreID  int
	FilePath string
	Accuracy float64
	Mods     osuapi.Mods
	MaxCombo int
	Misses   int
}

var tasks = make(chan oppaiTask, 500)

// Worker is a goroutine that calculates PP.
func Worker() {
	for task := range tasks {
		strmods := strings.TrimSpace(task.Mods.String())
		if strmods == "" {
			strmods = "nomod"
		}

		cmd := exec.Command(
			"./oppai",
			task.FilePath,
			strconv.FormatFloat(task.Accuracy, 'f', -1, 64)+"%",
			"+"+strings.ToLower(strmods),
			strconv.Itoa(task.MaxCombo)+"x",
			strconv.Itoa(task.Misses)+"m",
		)
		fmt.Println("=> oppai", cmd.Args)
		stdout := bytes.NewBuffer(nil)
		cmd.Stderr = os.Stderr
		cmd.Stdout = stdout
		err := cmd.Start()
		if err != nil {
			fmt.Println("=> pp calc failed:", err)
			setScore(task.ScoreID, failedCmd, 0)
			continue
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println("=> pp calc failed:", err)
			setScore(task.ScoreID, failedCmd, 0)
			continue
		}
		returned := string(stdout.Bytes())
		if len(returned) == 0 {
			setScore(task.ScoreID, failedInvalidReturnFormat, 0)
			continue
		}
		returnedF, err := strconv.ParseFloat(returned, 64)
		if err != nil {
			fmt.Println("=>", err)
			continue
		}
		err = setScore(task.ScoreID, completed, returnedF)
		if err != nil {
			fmt.Println("=>", err)
			continue
		}
		fmt.Println("=>", task.ScoreID, "calculated:", returnedF)
	}
}

func setScore(scoreID, calculatedStatus int, pp float64) error {
	_, err := db.Exec("UPDATE scores SET calculated = ?, total_pp = ? WHERE id = ?", calculatedStatus, pp, scoreID)
	return err
}
