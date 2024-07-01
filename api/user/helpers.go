package user

import (
	"errors"
	"github.com/google/uuid"
	"net/url"
	"strconv"
	"strings"
)

func validatePassportNumber(p string) (int, int, error) {
	raw := strings.Split(p, " ")
	var data []string

	for _, value := range raw {
		value = strings.TrimSpace(value)
		if value != "" {
			data = append(data, value)
		}
	}

	if len(data) != 2 {
		return 0, 0, errors.New("invalid passport number, use format: '1234 56780'")
	}

	if len(data[0]) != 4 {
		return 0, 0, errors.New("passport serie must be 4 characters long")
	}
	if len(data[1]) != 6 {
		return 0, 0, errors.New("passport number must be 6 characters long")
	}

	serie, err := strconv.Atoi(data[0])
	if err != nil {
		return 0, 0, errors.New("passport serie must contain only numbers")
	}

	number, err := strconv.Atoi(data[1])
	if err != nil {
		return 0, 0, errors.New("passport number must contain only numbers")
	}

	return serie, number, nil
}

func filtersMap(queryParams url.Values) map[string]interface{} {
	filters := map[string]interface{}{}

	passportSerie := queryParams.Get("passportSerie")
	if passportSerie != "" {
		ps, _ := strconv.Atoi(passportSerie)
		filters["passport_serie"] = ps
	}

	passportNumber := queryParams.Get("passportNumber")
	if passportNumber != "" {
		pn, _ := strconv.Atoi(passportNumber)
		filters["passport_number"] = pn
	}

	name := queryParams.Get("name")
	if name != "" {
		filters["name"] = name
	}

	surname := queryParams.Get("surname")
	if surname != "" {
		filters["surname"] = surname
	}

	patronymic := queryParams.Get("patronymic")
	if patronymic != "" {
		filters["patronymic"] = patronymic
	}

	address := queryParams.Get("address")
	if address != "" {
		filters["address"] = address
	}

	userId := queryParams.Get("userId")
	if userId != "" {
		uid, _ := uuid.Parse(userId)
		filters["user_id"] = uid
	}

	return filters
}
