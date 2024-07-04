package task

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type FullTask struct {
	gorm.Model `json:"-"`
	TaskId     uuid.UUID `json:"task_id" example:"00000000-0000-0000-0000-000000000000" extensions:"x-order=1"`
	OwnerId    uuid.UUID `json:"owner_id" example:"00000000-0000-0000-0000-000000000000" extensions:"x-order=2"`
	Title      string    `json:"title" example:"Title" extensions:"x-order=3"`
	Content    string    `json:"content" example:"Description" extensions:"x-order=4"`
	StartAt    time.Time `json:"start_at" example:"0001-01-01 00:00:00 +0000 UTC" extensions:"x-order=5"`
	FinishAt   time.Time `json:"end_at" example:"0001-01-01 00:00:00 +0000 UTC" extensions:"x-order=6"`
	Duration   int64     `json:"duration" example:"0" extensions:"x-order=7"`
}

type CreateTask struct {
	OwnerId uuid.UUID `json:"owner_id" extensions:"x-order=1"`
	Title   string    `json:"title" extensions:"x-order=2"`
	Content string    `json:"content" extensions:"x-order=3"`
}

type UpdateTask struct {
	TaskId  uuid.UUID `json:"-"`
	Title   string    `json:"title" extensions:"x-order=1"`
	Content string    `json:"content" extensions:"x-order=2"`
}

type OutputTask struct {
	Title    string `json:"title" extensions:"x-order=1"`
	Content  string `json:"content" extensions:"x-order=2"`
	Duration string `json:"duration" extensions:"x-order=3"`
}

type Summary struct {
	Name          string       `json:"name" extensions:"x-order=1"`
	Surname       string       `json:"surname" extensions:"x-order=2"`
	TasksDuration string       `json:"tasks_duration" extensions:"x-order=3"`
	Tasks         []OutputTask `json:"tasks" extensions:"x-order=4"`
}

func (f *FullTask) TableName() string {
	return "tasks"
}

func (f *FullTask) Create() error {
	err := DB.Create(f).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *FullTask) ReadOne() error {
	err := DB.Where("task_id = ?", f.TaskId).First(f).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *FullTask) ReadMany(filters map[string]time.Time) ([]FullTask, error) {
	var tasks []FullTask
	err := DB.
		Where("owner_id = ?", f.OwnerId).
		Where("finish_at BETWEEN ? and ?", filters["start_date"], filters["end_date"]).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (f *FullTask) UpdateFull() error {
	result := DB.Where("task_id = ?", f.TaskId).Updates(f)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (u *UpdateTask) UpdatePart() error {
	result := DB.Model(&FullTask{}).Where("task_id = ?", u.TaskId).Updates(u)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("404")
	}
	return nil
}

func (f *FullTask) Delete() error {
	result := DB.Where("task_id = ?", f.TaskId).Delete(f)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("404")
	}

	return nil
}
