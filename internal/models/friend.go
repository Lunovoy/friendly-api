package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID        uuid.UUID    `json:"id,omitempty" db:"id"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	DOB       sql.NullTime `json:"dob" db:"dob"`
	ImageID   uuid.UUID    `json:"image_id" db:"image_id"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
}

type UpdateFriendInput struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	DOB       *time.Time `json:"dob"`
	ImageID   *uuid.UUID `json:"image_id"`
}

type FriendWorkInfoTags struct {
	Friend   Friend   `json:"friend"`
	WorkInfo WorkInfo `json:"work_info"`
	Tags     []Tag    `json:"tags,omitempty"`
}

type UpdateFriendWorkInfoInput struct {
	Friend   *UpdateFriendInput   `json:"friend"`
	WorkInfo *UpdateWorkInfoInput `json:"work_info"`
	TagIDs   []AdditionTag        `json:"tag_ids"`
}

type FriendIDWorkInfoID struct {
	FriendID   uuid.UUID `json:"friend_id"`
	WorkInfoID uuid.UUID `json:"work_info_id"`
}

type FriendsTags struct {
	ID       uuid.UUID `json:"id,omitempty" db:"id"`
	FriendID uuid.UUID `json:"friend_id" db:"friend_id"`
	TagID    uuid.UUID `json:"tag_id" db:"tag_id"`
}

type FriendWithTags struct {
	Friend Friend `json:"friend"`
	Tags   []Tag  `json:"tags"`
}

type FriendID struct {
	FriendID uuid.UUID `json:"friend_id"`
}
