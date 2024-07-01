package user

import "net/http"

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/user", CreateUserHandler)
	router.HandleFunc("GET /api/v1/user/{uuid}", ReadUserByIDHandler)
	router.HandleFunc("GET /api/v1/user", ReadManyHandler)
	router.HandleFunc("PUT /api/v1/user/{uuid}", UpdateUserHandler)
	router.HandleFunc("DELETE /api/v1/user/{uuid}", DeleteUserHandler)
}
