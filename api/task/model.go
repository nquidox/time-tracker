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

func (t *Task) ReadMany() ([]Task, error) {
	var tasks []Task
	err := DB.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *Task) Update() error {
	result := DB.Where("task_id = ?", t.TaskId).Updates(t)

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

func (t *Task) Start() error {
	t.StartAt = time.Now()
	err := t.Update()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) GetStart() error {
	err := t.ReadOne()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Finish() error {
	err := t.Update()
	if err != nil {
		return err
	}
	return nil
}
