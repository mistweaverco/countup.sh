package main

import (
	"github.com/charmbracelet/log"
	"github.com/mistweaverco/timetrack.sh/internal/tui"
)

var VERSION string

func main() {
	log.Info("Starting timetrack.sh ⏰", "version", VERSION)
	tui.Start()
}