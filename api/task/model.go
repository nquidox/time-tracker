package task

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model `json:"-"`
	TaskId     uuid.UUID     `json:"task_id"`
	OwnerId    uuid.UUID     `json:"owner_id"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	StartAt    time.Time     `json:"start_at"`
	FinishAt   time.Time     `json:"end_at"`
	Duration   time.Duration `json:"duration"`
}

type UpdateTask struct {
	TaskId  uuid.UUID `json:"-"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type OutputTask struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Duration string `json:"duration"`
}

type Summary struct {
	Name          string       `json:"name"`
	Surname       string       `json:"surname"`
	TasksDuration string       `json:"tasks_duration"`
	Tasks         []OutputTask `json:"tasks"`
}

func (t *Task) Create() error {
	err := DB.Create(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) ReadOne() error {
	err := DB.Where("task_id = ?", t.TaskId).First(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) ReadMany(filters map[string]time.Time) ([]Task, error) {
	var tasks []Task
	err := DB.
		Where("owner_id = ?", t.OwnerId).
		Where("finish_at BETWEEN ? and ?", filters["start_date"], filters["end_date"]).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *Task) UpdateFull() error {
	result := DB.Where("task_id = ?", t.TaskId).Updates(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (t *UpdateTask) UpdatePart() error {
	result := DB.Model(&Task{}).Where("task_id = ?", t.TaskId).Updates(t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (t *Task) Delete() error {
	result := DB.Where("task_id = ?", t.TaskId).Delete(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
