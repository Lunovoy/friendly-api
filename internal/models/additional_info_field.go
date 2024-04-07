package models

import "github.com/google/uuid"

type AdditionalInfoField struct {
	ID     uuid.UUID `json:"id,omitempty" db:"id"`
	Title  string    `json:"title" db:"title"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}

type FriendsAdditionalInfoFields struct {
	ID                    uuid.UUID `json:"id,omitempty" db:"id"`
	FriendID              uuid.UUID `json:"friend_id" db:"friend_id"`
	AdditionalInfoFieldID uuid.UUID `json:"additional_info_field_id" db:"additional_info_field_id"`
}

type AdditionalInfoFieldText struct {
	ID                    uuid.UUID `json:"id,omitempty" db:"id"`
	Content               string    `json:"content" db:"content"`
	AdditionalInfoFieldID uuid.UUID `json:"additional_info_field_id" db:"additional_info_field_id"`
	FriendID              uuid.UUID `json:"friend_id" db:"friend_id"`
}
