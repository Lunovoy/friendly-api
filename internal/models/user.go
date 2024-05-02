package models

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	MiddleName string    `json:"middle_name" db:"middle_name"`
	TgUsername string    `json:"tg_username" db:"tg_username"`
	Mail       string    `json:"mail" binding:"required" db:"mail"`
	Password   string    `json:"password,omitempty" binding:"required" db:"password_hash"`
	Salt       string    `json:"salt,omitempty" db:"salt"`
}

type UserUpdate struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	MiddleName string    `json:"middle_name" db:"middle_name"`
	TgUsername string    `json:"tg_username" db:"tg_username"`
	Mail       string    `json:"mail" binding:"required" db:"mail"`
}
