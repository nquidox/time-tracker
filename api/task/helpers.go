package task

import (
	"errors"
	"github.com/google/uuid"
	"time_tracker/api/user"
)

func (t *Task) validateNewTask() error {
	p, err := uuid.Parse(t.OwnerId.String())
	if err != nil {
		return errors.New("incorrect user ID")
	}

	if p == uuid.Nil {
		return errors.New("owner ID is required")
	}

	var owner user.User
	result := DB.Where("user_id = ?", t.OwnerId).First(&owner)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task owner not found")
	}

	if len(t.Title) == 0 {
		return errors.New("title is required")
	}

	return nil
}
