package main

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func AppSetLogLevel(level string) log.Level {
	switch level {
	case "Trace":
		return log.TraceLevel
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warning":
		return log.WarnLevel
	case "Error":
		return log.ErrorLevel
	case "Fatal":
		return log.FatalLevel
	case "Panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func DBSetLogLevel(level string) logger.LogLevel {
	switch level {
	case "Info":
		return logger.Info
	case "Warning":
		return logger.Warn
	case "Error":
		return logger.Error
	case "Silent":
		return logger.Silent
	default:
		return logger.Info
	}
}
