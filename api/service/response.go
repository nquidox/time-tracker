package service

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ServerResponse(w http.ResponseWriter, dataInterface interface{}) {
	bytes, err := SerializeJSON(dataInterface)
	if err != nil {
		log.WithField("Serialize error", err).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.WithField("Response error", err).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
