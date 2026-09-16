package interviews

import (
	"errors"
	"time"
)

type ScheduleRequest struct {
	OrganizationID string
	ApplicationID  string
	CreatedBy      string
	StartsAt       time.Time
	EndsAt         time.Time
	Timezone       string
	Participants   []string
}

func ValidateSchedule(request ScheduleRequest, now time.Time) error {
	if request.OrganizationID == "" || request.ApplicationID == "" || request.CreatedBy == "" || request.Timezone == "" {
		return errors.New("organization, application, creator and timezone are required")
	}
	if !request.StartsAt.After(now) || !request.EndsAt.After(request.StartsAt) {
		return errors.New("interview time is invalid")
	}
	if len(request.Participants) == 0 {
		return errors.New("at least one participant is required")
	}
	return nil
}
