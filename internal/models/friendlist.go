package models

import (
	"github.com/google/uuid"
)

type Friendlist struct {
	ID          uuid.UUID `json:"id,omitempty" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
}

type FriendlistsTags struct {
	ID           uuid.UUID `json:"id,omitempty" db:"id"`
	FriendlistID uuid.UUID `json:"friendlist_id" db:"friendlist_id"`
	TagID        uuid.UUID `json:"tag_id" db:"tag_id"`
}

type FriendlistsFriends struct {
	ID           uuid.UUID `json:"id,omitempty" db:"id"`
	FriendlistID uuid.UUID `json:"friendlist_id" db:"friendlist_id"`
	FriendID     uuid.UUID `json:"friend_id" db:"friend_id"`
}

type FriendlistWithTags struct {
	Friendlist Friendlist `json:"friendlist"`
	Tags       []Tag      `json:"tags"`
}

type AdditionTagToFriendlist struct {
	TagID uuid.UUID `json:"tag_id" db:"tag_id" binding:"required"`
}

type AdditionFriendToFriendlist struct {
	FriendID uuid.UUID `json:"friend_id" db:"friend_id" binding:"required"`
}
