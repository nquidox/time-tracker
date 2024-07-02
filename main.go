package main

import (
	log "github.com/sirupsen/logrus"
	"time_tracker/api/task"
	"time_tracker/api/user"
	"time_tracker/config"
	"time_tracker/db"
)

func main() {
	c := config.New()

	log.SetLevel(AppSetLogLevel(c.Config.AppLogLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	user.ExternalAPIURL = c.Config.ExternalAPIURL

	DB := db.Connect(c, DBSetLogLevel(c.Config.DBLogLevel))
	user.Init(DB)
	task.Init(DB)

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
