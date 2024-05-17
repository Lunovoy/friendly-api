package models

import (
	"github.com/google/uuid"
)

type Reminder struct {
	ID                uuid.UUID `json:"id" db:"id"`
	MinutesUntilEvent int       `json:"minutes_until_event" db:"minutes_until_event"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	EventID           uuid.UUID `json:"event_id" db:"event_id"`
}

type ReminderUpdate struct {
	MinutesUntilEvent int `json:"minutes_until_event" db:"minutes_until_event"`
}

type ReminderWithIDUpdate struct {
	ID                *uuid.UUID `json:"id" db:"id"`
	MinutesUntilEvent *int       `json:"minutes_until_event" db:"minutes_until_event"`
}
