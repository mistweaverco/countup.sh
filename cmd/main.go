package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/countup.sh/internal/counter"
)

var VERSION string

func main() {
	log.Info("Starting countup.sh ⏰", "version", VERSION)

	if len(os.Args) < 2 {
		log.Error("No name for the timer provided 💀")
		os.Exit(1)
	}
	timerName := os.Args[1]

	log.Info("Timer started", "name", timerName)

	counter.Start(timerName)
}