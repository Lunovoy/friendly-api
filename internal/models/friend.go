package models

import (
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID         uuid.UUID `json:"id,omitempty" db:"id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	MiddleName string    `json:"middle_name" db:"middle_name"`
	DOB        time.Time `json:"dob" db:"dob"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
}

type FriendWorkInfo struct {
	Friend   Friend   `json:"friend" binding:"required"`
	WorkInfo WorkInfo `json:"work_info" binding:"required"`
}

type FriendIDWorkInfoID struct {
	FriendID   uuid.UUID `json:"friend_id"`
	WorkInfoID uuid.UUID `json:"work_info_id"`
}
