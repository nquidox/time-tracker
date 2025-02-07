package task

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(d *gorm.DB) {
	DB = d //passing DB global var
	err := DB.AutoMigrate(&FullTask{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Task model init success")
}
