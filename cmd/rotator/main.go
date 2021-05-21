package main

import (
	"flag"

	"github.com/Spazzy757/credentials-rotator/pkg/config"
	log "github.com/sirupsen/logrus"
)

var configHelpMessage = "The configuration file for credentials to rotate"

func main() {
	configFile := flag.String("config-file", "config.yaml", configHelpMessage)
	cfg := config.Config{}
	err := cfg.LoadConfig(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("config error")
	}
}
