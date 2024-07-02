package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time_tracker/config"
)

func Connect(c *config.Config, level logger.LogLevel) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Config.Host, c.Config.Port, c.Config.User, c.Config.Password, c.Config.Dbname, c.Config.Sslmode)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})

	if err != nil {
		log.Fatal(err)
	}
	return DB
}
