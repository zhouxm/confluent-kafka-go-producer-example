package main

import (
	"kafka-producer/cmd"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	//c := config.ParsedConfig.Parse()
	//log.Info(c)
	log.SetFormatter(&log.JSONFormatter{})
	// Set default log level; can be overwritten by configuration.
	log.SetLevel(log.InfoLevel)

	err := cmd.Execute()
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}
