package models

import "github.com/google/uuid"

type Tag struct {
	ID     uuid.UUID `json:"id,omitempty" db:"id"`
	Title  string    `json:"title" db:"title" binding:"required"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}
