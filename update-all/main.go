package main

import (
	"update-all/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
}

func main() {
	cmd.Execute()
}
