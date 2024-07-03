package task

import (
	"errors"
	"github.com/google/uuid"
	"net/url"
	"time"
	"time_tracker/api/user"
)

func (t *FullTask) validateNewTask() error {
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

func (t *UpdateTask) validateOnUpdate() error {
	if len(t.Title) == 0 {
		return errors.New("title is required")
	}
	return nil
}

func filtersMap(queryParams url.Values) map[string]time.Time {
	filters := map[string]time.Time{}
	layout := "2-1-2006"

	startDate := queryParams.Get("start_date")
	if startDate != "" {
		psd, err := time.Parse(layout, startDate)
		if err != nil {
			filters["start_date"] = time.Time{}
		}
		filters["start_date"] = psd
	} else {
		filters["start_date"] = time.Time{}
	}

	endDate := queryParams.Get("end_date")
	if endDate != "" {
		ped, err := time.Parse(layout, endDate)
		if err != nil {
			filters["end_date"] = time.Now()
		}
		filters["end_date"] = ped
	} else {
		filters["end_date"] = time.Now()
	}

	return filters
}
