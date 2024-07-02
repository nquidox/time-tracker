package user

import (
	"errors"
	"fmt"
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
	fmt.Printf("%s/info?passportSerie=%04d&passportNumber=%06d\n",
		ExternalAPIURL, passportSerie, passportNumber)

	response, err := http.Get(fmt.Sprintf("%s/info?passportSerie=%04d&passportNumber=%06d",
		ExternalAPIURL, passportSerie, passportNumber))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
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
		return errors.New("name field can't be empty")
	}

	if len(e.Surname) == 0 {
		return errors.New("surname field can't be empty")
	}

	if len(e.Address) == 0 {
		return errors.New("address field can't be empty")
	}

	return nil
}
