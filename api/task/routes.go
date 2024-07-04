package task

import "net/http"

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/task", CreateTaskHandler)
	router.HandleFunc("GET /api/v1/task/{uuid}", ReadOneTaskHandler)
	router.HandleFunc("GET /api/v1/tasks/{user_uuid}", ReadManyTaskHandler)
	router.HandleFunc("GET /api/v1/tasks/summary/{user_uuid}", SummaryHandler)
	router.HandleFunc("PUT /api/v1/task/{uuid}", UpdateTaskHandler)
	router.HandleFunc("DELETE /api/v1/task/{uuid}", DeleteTaskHandler)
	router.HandleFunc("GET /api/v1/task/start/{uuid}", StartTaskHandler)
	router.HandleFunc("GET /api/v1/task/finish/{uuid}", FinishTaskHandler)
}
