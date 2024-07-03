package task

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type FullTask struct {
	gorm.Model `json:"-"`
	TaskId     uuid.UUID     `json:"task_id" example:"00000000-0000-0000-0000-000000000000"`
	OwnerId    uuid.UUID     `json:"owner_id" example:"00000000-0000-0000-0000-000000000000"`
	Title      string        `json:"title" example:"Title"`
	Content    string        `json:"content" example:"Description"`
	StartAt    time.Time     `json:"start_at" example:"0001-01-01 00:00:00 +0000 UTC"`
	FinishAt   time.Time     `json:"end_at" example:"0001-01-01 00:00:00 +0000 UTC"`
	Duration   time.Duration `json:"duration" example:"0"`
}

type CreateTask struct {
	OwnerId uuid.UUID `json:"owner_id" example:"00000000-0000-0000-0000-000000000000"`
	Title   string    `json:"title" example:"New task title"`
	Content string    `json:"content" example:"Task description"`
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

func (t *FullTask) Create() error {
	err := DB.Create(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *FullTask) ReadOne() error {
	err := DB.Where("task_id = ?", t.TaskId).First(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *FullTask) ReadMany(filters map[string]time.Time) ([]FullTask, error) {
	var tasks []FullTask
	err := DB.
		Where("owner_id = ?", t.OwnerId).
		Where("finish_at BETWEEN ? and ?", filters["start_date"], filters["end_date"]).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *FullTask) UpdateFull() error {
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
	result := DB.Model(&FullTask{}).Where("task_id = ?", t.TaskId).Updates(t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (t *FullTask) Delete() error {
	result := DB.Where("task_id = ?", t.TaskId).Delete(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
