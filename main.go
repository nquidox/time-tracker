package main

import (
	"log"
	"time_tracker/api/user"
	"time_tracker/config"
	"time_tracker/db"
)

func main() {
	c := config.New()

	DB := db.Connect(c)
	user.Init(DB)

	server := NewApiServer(c.Config.HTTPHost, c.Config.HTTPPort)
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
