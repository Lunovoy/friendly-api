package models

import "github.com/google/uuid"

type User struct {
	ID                  uuid.UUID `json:"-" db:"id"`
	ImageID             uuid.UUID `json:"image_id" db:"image_id"`
	Username            string    `json:"username" db:"username"`
	FirstName           string    `json:"first_name" db:"first_name"`
	LastName            string    `json:"last_name" db:"last_name"`
	MiddleName          string    `json:"middle_name" db:"middle_name"`
	TgUsername          string    `json:"tg_username" db:"tg_username"`
	Mail                string    `json:"mail" binding:"required" db:"mail"`
	Password            string    `json:"password,omitempty" binding:"required" db:"password_hash"`
	Salt                string    `json:"-" db:"salt"`
	Country             string    `json:"country" db:"country"`
	City                string    `json:"city" db:"city"`
	Company             string    `json:"company" db:"company"`
	Profession          string    `json:"profession" db:"profession"`
	Position            string    `json:"position" db:"position"`
	Messenger           string    `json:"messenger" db:"messenger"`
	CommunicationMethod string    `json:"communication_method" db:"communication_method"`
	Nationality         string    `json:"nationality" db:"nationality"`
	Resident            bool      `json:"resident" db:"resident"`
	Language            string    `json:"language" db:"language"`
}

type UserUpdate struct {
	ID                  uuid.UUID `json:"-" db:"id"`
	ImageID             uuid.UUID `json:"image_id" db:"image_id"`
	Username            string    `json:"username" db:"username"`
	FirstName           string    `json:"first_name" db:"first_name"`
	LastName            string    `json:"last_name" db:"last_name"`
	MiddleName          string    `json:"middle_name" db:"middle_name"`
	TgUsername          string    `json:"tg_username" db:"tg_username"`
	Mail                string    `json:"mail" db:"mail"`
	Country             string    `json:"country" db:"country"`
	City                string    `json:"city" db:"city"`
	Company             string    `json:"company" db:"company"`
	Profession          string    `json:"profession" db:"profession"`
	Position            string    `json:"position" db:"position"`
	Messenger           string    `json:"messenger" db:"messenger"`
	CommunicationMethod string    `json:"communication_method" db:"communication_method"`
	Nationality         string    `json:"nationality" db:"nationality"`
	Resident            bool      `json:"resident" db:"resident"`
	Language            string    `json:"language" db:"language"`
}
