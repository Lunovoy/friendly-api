package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user models.User) (uuid.UUID, error) {
	var id uuid.UUID

	query := fmt.Sprintf("INSERT INTO \"%s\" (username, first_name, last_name, middle_name, tg_username, mail, password_hash, salt, country, city, company, profession, position, messenger, communication_method, nationality, language, resident) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id", userTable)

	row := r.db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.MiddleName, user.TgUsername, user.Mail, user.Password, user.Salt, user.Country, user.City, user.Company, user.Profession, user.Position, user.Messenger, user.CommunicationMethod, user.Nationality, user.Language, user.Resident)

	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByMail(mail string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE mail = $1", userTable)

	err := r.db.Get(&user, query, mail)
	return user, err
}

func (r *AuthPostgres) GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE id = $1", userTable)

	err := r.db.Get(&user, query, id)
	return user, err
}
