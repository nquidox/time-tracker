package user

import (
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d //passing DB global var
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}
