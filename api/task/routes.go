package task

import "net/http"

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/task", CreateTaskHandler)
	router.HandleFunc("GET /api/v1/task/{uuid}", ReadOneTaskHandler)

	//вместо проверки авторизации просто запрашиваем id пользователя
	router.HandleFunc("GET /api/v1/task/summary/{user_uuid}", ReadManyTaskHandler)

	router.HandleFunc("PUT /api/v1/task/{uuid}", UpdateTaskHandler)
	router.HandleFunc("DELETE /api/v1/task/{uuid}", DeleteTaskHandler)
	router.HandleFunc("GET /api/v1/task/start/{uuid}", StartTaskHandler)
	router.HandleFunc("GET /api/v1/task/stop/{uuid}", FinishTaskHandler)
}

// для всех роутов предполагаем, что они доступны только пользователям
// с соответствующими uuid или пользователями с особым статусом (админ)
