package main

import (
	"flag"

	"github.com/Spazzy757/credentials-rotator/pkg/config"
	"github.com/Spazzy757/credentials-rotator/pkg/handlers"
	log "github.com/sirupsen/logrus"
)

var configHelpMessage = "The configuration file for credentials to rotate"

func main() {
	// Flags
	configFile := flag.String("config-file", "config.yaml", configHelpMessage)
	flag.Parse()

	cfg := config.Config{}
	err := cfg.LoadConfig(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("config error")
	}
	err = handlers.ConfigHandler(&cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("runtime error")
	}
	log.WithFields(log.Fields{
		"count": len(cfg.Credentials),
	}).Info("success")
}
