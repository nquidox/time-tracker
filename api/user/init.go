package user

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d //passing DB global var
	err := DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("User model init success")
}
