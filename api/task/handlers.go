package task

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"time"
	"time_tracker/api/service"
	"time_tracker/api/user"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	tsk := Task{TaskId: uuid.New()}

	err = service.DeserializeJSON(data, &tsk)
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}

	err = tsk.validateNewTask()
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	err = tsk.Create()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "Task created successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    tsk.TaskId,
	})
	log.Info(msg)
}

func ReadOneTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}
	err = tsk.ReadOne()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, tsk)
	log.Info("Read one successfully")
}

func ReadManyTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	queryParams := r.URL.Query()
	filters := filtersMap(queryParams)

	userId, err := uuid.Parse(r.PathValue("user_uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{OwnerId: userId}

	// для простоты выборка задач для расчета трудозатрат будет производиться по полю finish_at
	// если требуется также учитывать промежуточное состояние, когда задача начата,
	// но еще не закончена на текущий момент, то желателен механизм паузы
	// перенести потом в ридми
	tasks, err := tsk.ReadMany(filters)
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Duration > tasks[j].Duration
	})

	sumDuration := time.Duration(0)
	for _, t := range tasks {
		sumDuration += t.Duration
	}

	var outputList []OutputTask
	for _, tsk := range tasks {
		outputList = append(outputList, OutputTask{
			Title:   tsk.Title,
			Content: tsk.Content,
			Duration: fmt.Sprintf("%02d:%02d:%02d",
				int(tsk.Duration.Hours()),
				int(tsk.Duration.Minutes())%60,
				int(tsk.Duration.Seconds())%60),
		})
	}

	usr := user.User{UserId: userId}
	err = usr.ReadOne()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	response := Summary{
		Name:    usr.Name,
		Surname: usr.Surname,
		TasksDuration: fmt.Sprintf("%02d:%02d:%02d",
			int(sumDuration.Hours()),
			int(sumDuration.Minutes())%60,
			int(sumDuration.Seconds())%60),
		Tasks: outputList,
	}

	service.ServerResponse(w, response)
	log.Info("Read many successfully")
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}
	defer r.Body.Close()

	tsk := UpdateTask{TaskId: taskId}

	err = service.DeserializeJSON(data, &tsk)
	if err != nil {
		e.DeserializeError(err)
		service.ServerResponse(w, e)
		return
	}

	err = tsk.validateOnUpdate()
	if err != nil {
		e.ValidationError(err)
		service.ServerResponse(w, e)
		return
	}

	err = tsk.UpdatePart()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "Task updated successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    "",
	})
	log.Info(msg)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}

	err = tsk.Delete()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "Task deleted successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    "",
	})
	log.Info(msg)
}

func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}
	err = tsk.ReadOne()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	if !tsk.StartAt.IsZero() {
		e.TaskIsAlreadyStartedError()
		service.ServerResponse(w, e)
		return
	}

	tsk.StartAt = time.Now()

	err = tsk.UpdateFull()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "Task started successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    fmt.Sprintf("Started at: %s", tsk.StartAt.Format("15:04:05 02-01-2006")),
	})
	log.Info(msg)
}

func FinishTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	log.Info(r.Method, " ", r.URL.Path, " ", r.RemoteAddr, " ", r.UserAgent())

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}

	err = tsk.ReadOne()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	fmt.Println(tsk.StartAt, tsk.StartAt.IsZero())

	if tsk.StartAt.IsZero() {
		e.TaskNotStartedError()
		service.ServerResponse(w, e)
		return
	}

	if !tsk.FinishAt.IsZero() {
		e.TaskIsAlreadyFinishedError()
		service.ServerResponse(w, e)
		return
	}

	tsk.FinishAt = time.Now()
	tsk.Duration = tsk.FinishAt.Sub(tsk.StartAt)

	err = tsk.UpdateFull()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	msg := "Task finished successfully"

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data: fmt.Sprintf("Finished at: %s, Duration: %02d:%02d:%02d",
			tsk.FinishAt.Format("15:04:05 02-01-2006"),
			int(tsk.Duration.Hours()),
			int(tsk.Duration.Minutes())%60,
			int(tsk.Duration.Seconds())%60,
		),
	})
	log.Info(msg)
}
