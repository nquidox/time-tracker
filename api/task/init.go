package task

import (
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d //passing DB global var
	err := DB.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal(err)
	}
}
