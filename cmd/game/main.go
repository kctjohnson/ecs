package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"ecs/internal/game"
	"ecs/internal/ui"
)

func main() {
	debug := false
	var logger *log.Logger
	if debug {
		logFileName := fmt.Sprintf("log_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = log.New(logFile, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	game := game.NewGame(logger)
	ui.RunGame(game, logger)
}
