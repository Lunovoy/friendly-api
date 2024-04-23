package models

import (
	"github.com/google/uuid"
)

type Friendlist struct {
	ID          uuid.UUID `json:"id,omitempty" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	ImageID     uuid.UUID `json:"image_id" db:"image_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
}

type UpdateFriendlist struct {
	Title       *string    `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	ImageID     *uuid.UUID `json:"image_id" db:"image_id"`
	UserID      *uuid.UUID `json:"user_id" db:"user_id"`
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

type FriendlistFull struct {
	FriendlistWithTags FriendlistWithTags   `json:"friendlist"`
	FriendsWithTags    []FriendWorkInfoTags `json:"friends"`
}

type FriendlistWithFriends struct {
	Friendlist Friendlist `json:"friendlist"`
	Friends    []Friend   `json:"friends"`
}

type AdditionFriendToFriendlist struct {
	FriendID uuid.UUID `json:"friend_id" db:"friend_id" binding:"required"`
}
