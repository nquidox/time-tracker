package task

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
	"time_tracker/api/service"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

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

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "Task created successfully",
		Data:    tsk.TaskId,
	})
}

func ReadOneTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

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
}

func ReadManyTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e service.ErrorResponse
	var task Task

	tasks, err := task.ReadMany()
	if err != nil {
		e.ReadBodyError(err)
		service.ServerResponse(w, e)
		return
	}

	//TODO add user auth and sort response by duration
	service.ServerResponse(w, tasks)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

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

	tsk := Task{TaskId: taskId}

	err = service.DeserializeJSON(data, &tsk)
	if err != nil {
		e.DeserializeError(err)
		service.ServerResponse(w, e)
		return
	}

	err = tsk.Update()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "Task updated successfully",
		Data:    "",
	})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

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

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "Task deleted successfully",
		Data:    "",
	})
}

func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}

	//TODO check if already started
	err = tsk.Start()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "Task started successfully",
		Data:    fmt.Sprintf("Started at: %s", tsk.StartAt.Format("15:04:05 02-01-2006")),
	})
}

func FinishTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e service.ErrorResponse

	taskId, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		e.UuidParseError(err)
		service.ServerResponse(w, e)
		return
	}

	tsk := Task{TaskId: taskId}
	zeroTime := time.Time{}

	err = tsk.GetStart()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}
	//FIXME fix zero time check
	if tsk.StartAt == zeroTime {
		service.ServerResponse(w, service.ErrorResponse{
			Code:    400,
			Message: "Task not started",
		})
		return
	}

	if tsk.FinishAt != zeroTime {
		service.ServerResponse(w, service.ErrorResponse{
			Code:    400,
			Message: "Task is already finished",
		})
		return
	}

	tsk.FinishAt = time.Now()
	tsk.Duration = tsk.FinishAt.Sub(tsk.StartAt)

	err = tsk.Finish()
	if err != nil {
		e.DBError(err)
		service.ServerResponse(w, e)
		return
	}

	service.ServerResponse(w, service.OkResponse{
		Code:    http.StatusOK,
		Message: "Task finished successfully",
		Data:    fmt.Sprintf("Finished at: %s, Duration: %d", tsk.FinishAt.Format("15:04:05 02-01-2006"), tsk.Duration),
	})
}
