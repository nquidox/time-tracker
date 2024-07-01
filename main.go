package main

import (
	"fmt"
	"time_tracker/config"
	"time_tracker/db"
)

func main() {
	c := config.New()
	fmt.Println(c)

	DB := db.Connect(c)
	fmt.Println(DB)
}
