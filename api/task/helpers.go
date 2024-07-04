package task

import (
	"errors"
	"github.com/google/uuid"
	"net/url"
	"time"
	"time_tracker/api/user"
)

func (f *FullTask) validateNewTask() error {
	p, err := uuid.Parse(f.OwnerId.String())
	if err != nil {
		return errors.New("incorrect user ID")
	}

	if p == uuid.Nil {
		return errors.New("owner ID is required")
	}

	err = validateOwner(f.OwnerId)
	if err != nil {
		return err
	}

	if len(f.Title) == 0 {
		return errors.New("title is required")
	}

	return nil
}

func (u *UpdateTask) validateOnUpdate() error {
	if len(u.Title) == 0 {
		return errors.New("title can't be ommited or be blank")
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

func validateOwner(id uuid.UUID) error {
	var owner user.FullUser
	err := DB.Where("user_id = ?", id).First(&owner).Error
	if err != nil {
		return err
	}
	return nil
}
