package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	userTable                        = "user"
	tagTable                         = "tag"
	friendlistTable                  = "friendlist"
	friendlistsTagsTable             = "friendlists_tags"
	friendlistsFriendsTable          = "friendlists_friends"
	friendTable                      = "friend"
	workInfoTable                    = "work_info"
	additionalInfoFieldTable         = "additional_info_field"
	friendsAdditionalInfoFieldsTable = "friends_additional_info_fields"
	additionalInfoFieldTextTable     = "additional_info_field_text"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}