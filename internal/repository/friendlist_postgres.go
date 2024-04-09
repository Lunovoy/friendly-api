package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type FriendlistPostgres struct {
	db *sqlx.DB
}

func NewFriendlistPostgres(db *sqlx.DB) *FriendlistPostgres {
	return &FriendlistPostgres{
		db: db,
	}
}

func (r *FriendlistPostgres) Create(userID uuid.UUID, friendlist models.Friendlist) (uuid.UUID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var friendlistID uuid.UUID
	query := fmt.Sprintf("INSERT INTO \"%s\" (title, description, user_id) VALUES ($1, $2, $3) RETURNING id", friendlistTable)

	row := tx.QueryRow(query, friendlist.Title, friendlist.Description, userID)
	if err := row.Scan(&friendlistID); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return friendlistID, err
}

func (r *FriendlistPostgres) GetAll(userID uuid.UUID) ([]models.Friendlist, error) {
	var friendlists []models.Friendlist

	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s where user_id = $1", friendlistTable)

	err := r.db.Select(&friendlists, query, userID)

	return friendlists, err

}

func (r *FriendlistPostgres) GetByID(userID, friendlistID uuid.UUID) (models.Friendlist, error) {
	var friendlist models.Friendlist

	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE id = $1 AND user_id = $2", friendlistTable)

	err := r.db.Get(&friendlist, query, friendlistID, userID)

	return friendlist, err
}

func (r *FriendlistPostgres) GetAllWithTags(userID uuid.UUID) ([]models.FriendlistWithTags, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var friendlists []models.Friendlist

	queryFriendlists := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE user_id = $1", friendlistTable)

	if err := tx.Select(&friendlists, queryFriendlists, userID); err != nil {
		return nil, err
	}

	queryTags := fmt.Sprintf(`SELECT t.id, t.title, t.user_id 
							FROM %s t
							INNER JOIN %s ft ON ft.tag_id = t.id 
							WHERE ft.friendlist_id = $1`, tagTable, friendlistsTagsTable)

	stmt, err := tx.Preparex(queryTags)
	if err != nil {
		return nil, err
	}

	var friendlistsWithTags []models.FriendlistWithTags

	var tags []models.Tag
	for _, friendlist := range friendlists {

		if err := stmt.Select(&tags, friendlist.ID); err != nil {
			return nil, err
		}
		friendlistsWithTags = append(friendlistsWithTags, models.FriendlistWithTags{
			Friendlist: friendlist,
			Tags:       tags,
		})
		tags = nil
	}

	return friendlistsWithTags, err

}

func (r *FriendlistPostgres) GetByIDWithTags(userID, friendlistID uuid.UUID) (models.FriendlistWithTags, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendlistWithTags{}, err
	}
	defer tx.Rollback()

	var friendlist models.Friendlist

	queryFriendlist := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE id = $1 AND user_id = $2", friendlistTable)

	if err := tx.Get(&friendlist, queryFriendlist, friendlistID, userID); err != nil {
		return models.FriendlistWithTags{}, err
	}

	queryTags := fmt.Sprintf(`SELECT t.id, t.title, t.user_id 
							FROM %s t
							INNER JOIN %s ft ON ft.tag_id = t.id 
							WHERE ft.friendlist_id = $1`, tagTable, friendlistsTagsTable)

	var tags []models.Tag
	if err := tx.Select(&tags, queryTags, friendlist.ID); err != nil {
		return models.FriendlistWithTags{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.FriendlistWithTags{}, err
	}

	friendlistWithTags := models.FriendlistWithTags{
		Friendlist: friendlist,
		Tags:       tags,
	}

	return friendlistWithTags, err
}

func (r *FriendlistPostgres) Update(userID, friendlistID uuid.UUID, friendlist models.Friendlist) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2 WHERE id = $3 AND user_id = $4", friendlistTable)

	_, err := r.db.Exec(query, friendlist.Title, friendlist.Description, friendlistID, userID)

	return err
}

func (r *FriendlistPostgres) DeleteByID(userID, friendlistID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s where id=$1 AND user_id=$2", friendlistTable)

	_, err := r.db.Exec(query, friendlistID, userID)
	if err != nil {
		return err
	}

	return err
}

func (r *FriendlistPostgres) AddTagToFriendlist(friendlistID, tagID uuid.UUID) error {

	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friendlist_id = $1 AND tag_id = $2)", friendlistsTagsTable)

	if err := r.db.Get(&exists, queryCheck, friendlistID, tagID); err != nil {
		return err
	}
	if exists {
		return errors.New("tag already exists in this friendlist")
	}

	queryAdd := fmt.Sprintf("INSERT INTO \"%s\" (friendlist_id, tag_id) VALUES ($1, $2)", friendlistsTagsTable)

	_, err := r.db.Exec(queryAdd, friendlistID, tagID)

	return err
}

func (r *FriendlistPostgres) AddFriendToFriendlist(friendlistID, friendID uuid.UUID) error {
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friendlist_id = $1 AND friend_id = $2)", friendlistsFriendsTable)

	if err := r.db.Get(&exists, queryCheck, friendlistID, friendID); err != nil {
		return err
	}
	if exists {
		return errors.New("friend already exists in friendlist")
	}

	queryAdd := fmt.Sprintf("INSERT INTO \"%s\" (friendlist_id, friend_id) VALUES ($1, $2)", friendlistsFriendsTable)

	_, err := r.db.Exec(queryAdd, friendlistID, friendID)

	return err
}

func (r *FriendlistPostgres) DeleteTagFromFriendlist(friendlistID, tagID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s where friendlist_id=$1 AND tag_id=$2", friendlistsTagsTable)

	_, err := r.db.Exec(query, friendlistID, tagID)
	if err != nil {
		return err
	}

	return err
}

func (r *FriendlistPostgres) DeleteFriendFromFriendlist(friendlistID, friendID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s where friendlist_id=$1 AND friend_id=$2", friendlistsFriendsTable)

	_, err := r.db.Exec(query, friendlistID, friendID)
	if err != nil {
		return err
	}

	return err
}
