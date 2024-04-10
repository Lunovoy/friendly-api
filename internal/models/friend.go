package models

import (
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	DOB       time.Time `json:"dob" db:"dob"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
}
type UpdateFriendInput struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	DOB       *time.Time `json:"dob"`
}

type FriendWorkInfo struct {
	Friend   Friend   `json:"friend" binding:"required"`
	WorkInfo WorkInfo `json:"work_info" binding:"required"`
}
type UpdateFriendWorkInfoInput struct {
	Friend   *UpdateFriendInput   `json:"friend"`
	WorkInfo *UpdateWorkInfoInput `json:"work_info"`
}

type FriendIDWorkInfoID struct {
	FriendID   uuid.UUID `json:"friend_id"`
	WorkInfoID uuid.UUID `json:"work_info_id"`
}
