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

	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE id = $1", userTable)

	err := r.db.Get(&user, query, userID)

	return user, err
}

func (r *UserPostgres) Update(user models.UserUpdate, userID uuid.UUID) error {
	query := fmt.Sprintf("UPDATE \"%s\" SET username = $1, first_name = $2, last_name = $3, middle_name = $4, tg_username = $5, mail = $6, country = $7, city = $8, company = $9, profession = $10, position = $11, messenger = $12, communication_method = $13, nationality = $14, language = $15, resident = $16, image_id = $17 WHERE id = $18", userTable)

	_, err := r.db.Exec(query, user.Username, user.FirstName, user.LastName, user.MiddleName, user.TgUsername, user.Mail, user.Country, user.City, user.Company, user.Profession, user.Position, user.Messenger, user.CommunicationMethod, user.Nationality, user.Language, user.Resident, user.ImageID, userID)

	return err
}
