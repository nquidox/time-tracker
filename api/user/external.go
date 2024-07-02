package user

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time_tracker/api/service"
)

type ExternalUser struct {
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Surname    string `json:"surname"`
	Address    string `json:"address"`
}

var ExternalAPIURL string

func (e *ExternalUser) GetExternalData(passportSerie, passportNumber int) error {
	link := fmt.Sprintf("%s/info?passportSerie=%04d&passportNumber=%06d",
		ExternalAPIURL, passportSerie, passportNumber)

	log.WithField("Link", link).Debug("Retrieving external user data")

	response, err := http.Get(link)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("External user data retrieval error: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = service.DeserializeJSON(data, e)
	if err != nil {
		return err
	}

	return nil
}

func (e *ExternalUser) ValidateRequiredFields() error {
	if len(e.Name) == 0 {
		msg := "name field can't be empty"
		log.Error(msg)
		return errors.New(msg)
	}

	if len(e.Surname) == 0 {
		msg := "surname field can't be empty"
		log.Error(msg)
		return errors.New(msg)
	}

	if len(e.Address) == 0 {
		msg := "address field can't be empty"
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}
