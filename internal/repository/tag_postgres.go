package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type TagPostgres struct {
	db *sqlx.DB
}

func NewTagPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{
		db: db,
	}
}

func (r *TagPostgres) Create(userID uuid.UUID, tag models.Tag) (uuid.UUID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE title = $1 AND user_id = $2)", tagTable)

	if err := r.db.Get(&exists, queryCheck, tag.Title, userID); err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("tag already exists")
	}

	var tagID uuid.UUID
	query := fmt.Sprintf("INSERT INTO \"%s\" (title, user_id) VALUES ($1, $2) RETURNING id", tagTable)

	row := tx.QueryRow(query, tag.Title, userID)
	if err := row.Scan(&tagID); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return tagID, nil
}

func (r *TagPostgres) GetAll(userID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag

	query := fmt.Sprintf("SELECT id, title, user_id FROM %s where user_id = $1", tagTable)

	err := r.db.Select(&tags, query, userID)

	return tags, err

}

func (r *TagPostgres) GetByID(userID, tagID uuid.UUID) (models.Tag, error) {
	var tag models.Tag

	query := fmt.Sprintf("SELECT id, title, user_id FROM %s WHERE id = $1 AND user_id = $2", tagTable)

	err := r.db.Get(&tag, query, tagID, userID)

	return tag, err
}

func (r *TagPostgres) Update(userID, tagID uuid.UUID, tag models.Tag) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1 WHERE id = $2 AND user_id = $3", tagTable)

	_, err := r.db.Exec(query, tag.Title, tagID, userID)

	return err
}

func (r *TagPostgres) DeleteByID(userID, tagID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND user_id=$2", tagTable)

	_, err := r.db.Exec(query, tagID, userID)

	return err
}
