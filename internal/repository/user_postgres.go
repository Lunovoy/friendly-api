package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (r *UserPostgres) GetByID(userID uuid.UUID) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id, username, first_name, last_name, middle_name, tg_username, mail FROM \"%s\" WHERE id = $1", userTable)

	err := r.db.Get(&user, query, userID)

	return user, err
}

func (r *UserPostgres) Update(user models.UserUpdate, userID uuid.UUID) error {
	query := fmt.Sprintf("UPDATE \"%s\" SET username = $1, first_name = $2, last_name = $3, middle_name = $4, tg_username = $5, mail = $6 WHERE id = $7", userTable)

	_, err := r.db.Exec(query, user.Username, user.FirstName, user.LastName, user.MiddleName, user.TgUsername, user.Mail, userID)

	return err
}
