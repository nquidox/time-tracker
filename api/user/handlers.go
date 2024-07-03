package user

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time_tracker/api/service"
)

// CreateUserHandler godoc
//
//	@Summary		Create user
//	@Description	Create user by passport serie and number
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			New	user		body	NewUser	true	"Provide passport serie and number in format '1234 567890'"
//	@Success		200	{object}	FullUser
//	@Failure		400	{object}	service.ErrorResponse
//	@Failure		500	{object}	service.ErrorResponse
//	@Router			/user [post]
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

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

	if exists(serie, number) != uuid.Nil {
		e.DBExists()
		service.ServerResponse(w, e)
		log.Error("User already exists")
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

	usr := FullUser{
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

	msg := "User created successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    usr.UserId,
	})
	log.WithField("Full name", fmt.Sprintf("%s %s %s", usr.Name, usr.Patronymic, usr.Surname)).
		Info(msg)
}

// ReadUserByIDHandler godoc
//
//	@Summary		Get user
//	@Description	Get user by UUID
//	@Tags			User
//	@Produce		json
//	@Param			uuid	path		string	true	"Provide user's uuid"
//	@Success		200		{object}	FullUser
//	@Failure		400		{object}	service.ErrorResponse
//	@Failure		404		{object}	service.ErrorResponse
//	@Failure		500		{object}	service.ErrorResponse
//	@Router			/user/{uuid} [get]
func ReadUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := FullUser{UserId: userId}
	err = usr.ReadOne()
	if err != nil {
		if err.Error() == "record not found" {
			e.Error404()
		} else {
			e.DBError(err)
		}
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, usr)
	log.Info("User read successfully")
}

// ReadManyHandler godoc
//
//	@Summary		Get all users
//	@Description	Get all users with filters and pagination
//	@Tags			User
//	@Produce		json
//	@Param			passportSerie	query		int		false	"Passport serie"
//	@Param			passportNumber	query		int		false	"Passport number"
//	@Param			name			query		string	false	"Name"
//	@Param			surname			query		string	false	"Surname"
//	@Param			patronymic		query		string	false	"Patronymic"
//	@Param			address			query		string	false	"Address"
//	@Param			userId			query		string	false	"User UUID"
//	@Param			page			query		int		false	"Page number"
//	@Param			perPage			query		int		false	"Records per page"
//
//	@Success		200				{object}	FullUser
//	@Failure		400				{object}	service.ErrorResponse
//	@Failure		404				{object}	service.ErrorResponse
//	@Failure		500				{object}	service.ErrorResponse
//	@Router			/user [get]
func ReadManyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	queryParams := r.URL.Query()
	filters := filtersMap(queryParams)
	params := paginationParams(queryParams)

	var usr FullUser
	users, err := usr.ReadMany(filters, params)
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	if len(users) == 0 {
		e.Error404()
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, users)
	log.WithFields(log.Fields{
		"page":     params["page"],
		"per_page": params["per_page"],
	}).Info("Users read successfully")
}

// UpdateUserHandler godoc
//
//	@Summary		Update user
//	@Description	Update user by UUID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true		"Provide user's uuid"
//	@Param			User	data		body	FullUser	true	"Passport serie and number are required. Partial update possible, empty fields will be ignored."
//	@Success		200		{object}	service.OkResponse
//	@Failure		400		{object}	service.ErrorResponse
//	@Failure		404		{object}	service.ErrorResponse
//	@Failure		500		{object}	service.ErrorResponse
//	@Router			/user/{uuid} [put]
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := UpdateUser{UserId: userId}

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

	pn := fmt.Sprintf("%s %s", strconv.Itoa(usr.PassportSerie), strconv.Itoa(usr.PassportNumber))
	_, _, err = validatePassportNumber(pn)
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	id := exists(usr.PassportSerie, usr.PassportNumber)
	if id != usr.UserId && id != uuid.Nil {
		log.Error(usr.UserId)
		e.DBPassportExists()
		service.ServerResponse(w, e)
		return
	}

	err = usr.Update()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "User updated successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    204,
		Message: msg,
		Data:    "",
	})
	log.Info(msg)
}

// DeleteUserHandler godoc
//
//	@Summary		Delete user
//	@Description	Delete user by UUID
//	@Tags			User
//	@Produce		json
//	@Param			uuid	path		string	true	"Provide user's uuid"
//	@Success		200		{object}	service.OkResponse
//	@Failure		400		{object}	service.ErrorResponse
//	@Failure		404		{object}	service.ErrorResponse
//	@Failure		500		{object}	service.ErrorResponse
//	@Router			/user/{uuid} [delete]
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	userId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	usr := FullUser{UserId: userId}

	err = usr.Delete()
	if err != nil {
		e.Error404()
		service.ServerResponse(w, e)
		return
	}

	msg := "User deleted successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    204,
		Message: msg,
		Data:    "",
	})
	log.Info(msg)
}
