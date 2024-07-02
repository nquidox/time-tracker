package user

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time_tracker/api/service"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	data, err := io.ReadAll(r.Body)
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}
	defer r.Body.Close()

	newUsr := NewUser{}
	err = service.DeserializeJSON(data, &newUsr)
	if err != nil {
		e.DeserializeError(err)
		service.ServerResponse(w, e)
		return
	}

	serie, number, err := validatePassportNumber(newUsr.PassportNumber)
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	var extUser ExternalUser
	err = extUser.GetExternalData(serie, number)
	if err != nil {
		e.ExternalAPIError(err)
		service.ServerResponse(w, e)
		return
	}

	err = extUser.ValidateRequiredFields()
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := User{
		PassportSerie:  serie,
		PassportNumber: number,
		Name:           extUser.Name,
		Surname:        extUser.Surname,
		Patronymic:     extUser.Patronymic,
		Address:        extUser.Address,
		UserId:         uuid.New(),
	}

	err = usr.Create()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "User created successfully",
		Data:    usr.UserId,
	})
}

func ReadUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := User{UserId: userId}
	err = usr.ReadOne()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		e.SerializeError(err)
		service.ServerResponse(w, e)
		return
	}
}

func ReadManyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	queryParams := r.URL.Query()
	filters := filtersMap(queryParams)
	params := paginationParams(queryParams)

	var usr User
	users, err := usr.ReadMany(filters, params)
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, users)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := User{UserId: userId}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}
	defer r.Body.Close()

	err = service.DeserializeJSON(data, &usr)
	if err != nil {
		e.DeserializeError(err)
		service.ServerResponse(w, e)
		return
	}

	err = usr.Update()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    204,
		Message: "User updated successfully",
		Data:    "",
	})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := User{UserId: userId}

	err = usr.Delete()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    204,
		Message: "User updated successfully",
		Data:    "",
	})
}
