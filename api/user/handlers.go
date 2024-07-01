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

	newUsr := NewUser{}
	err = service.DeserializeJSON(data, &newUsr)
	if err != nil {
		e.DeserializeError(err)
		service.ServerResponse(w, e)
		return
	}

	s, n, err := validatePassportNumber(newUsr.PassportNumber)
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := User{
		PassportSerie:  s,
		PassportNumber: n,
		Name:           "request from external api",
		Surname:        "request from external api",
		Patronymic:     "request from external api",
		Address:        "request from external api",
		UserId:         uuid.New(),
	}

	err = usr.Create()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    200,
		Message: "Ok",
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

	var usr User
	users, err := usr.ReadMany(filters)
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		e.SerializeError(err)
		service.ServerResponse(w, e)
		return
	}
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
