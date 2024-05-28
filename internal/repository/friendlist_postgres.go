package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
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

func (r *FriendlistPostgres) Create(userID uuid.UUID, friendlist models.UpdateFriendlist) (uuid.UUID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	friendlistFields := []string{"title", "user_id"}
	friendlistValues := []any{*friendlist.Title, userID}
	builderFriendlist := sqlbuilder.NewInsertBuilder()
	builderFriendlist.SetFlavor(sqlbuilder.PostgreSQL)
	builderFriendlist.InsertInto(friendlistTable)

	if friendlist.Description != nil {
		friendlistFields = append(friendlistFields, "description")
		friendlistValues = append(friendlistValues, *friendlist.Description)
	}
	if friendlist.Color != nil {
		friendlistFields = append(friendlistFields, "color")
		friendlistValues = append(friendlistValues, *friendlist.Color)
	}
	if friendlist.ImageID != nil {
		friendlistFields = append(friendlistFields, "image_id")
		friendlistValues = append(friendlistValues, *friendlist.ImageID)
	}

	builderFriendlist.Cols(friendlistFields...).Values(friendlistValues...)

	queryFriendlist, args := builderFriendlist.Build()
	queryFriendlist += " RETURNING id;"

	var friendlistID uuid.UUID

	row := tx.QueryRow(queryFriendlist, args...)
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

	query := fmt.Sprintf("SELECT * FROM %s where user_id = $1", friendlistTable)

	err := r.db.Select(&friendlists, query, userID)

	return friendlists, err

}

func (r *FriendlistPostgres) GetByID(userID, friendlistID uuid.UUID) (models.Friendlist, error) {
	var friendlist models.Friendlist

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", friendlistTable)

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

	queryFriendlists := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", friendlistTable)

	if err := tx.Select(&friendlists, queryFriendlists, userID); err != nil {
		return nil, err
	}

	queryTags := fmt.Sprintf(`SELECT t.* 
							FROM %s t
							INNER JOIN %s ft ON ft.tag_id = t.id 
							WHERE ft.friendlist_id = $1`, tagTable, friendlistsTagsTable)

	tagStmt, err := tx.Preparex(queryTags)
	if err != nil {
		return nil, err
	}

	var friendlistsWithTags []models.FriendlistWithTags

	var tags []models.Tag
	for _, friendlist := range friendlists {

		if err := tagStmt.Select(&tags, friendlist.ID); err != nil {
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

	queryFriendlist := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", friendlistTable)

	if err := tx.Get(&friendlist, queryFriendlist, friendlistID, userID); err != nil {
		return models.FriendlistWithTags{}, err
	}

	queryTags := fmt.Sprintf(`SELECT t.* 
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

func (r *FriendlistPostgres) GetAllWithFriends(userID uuid.UUID) ([]models.FriendlistWithFriends, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var friendlists []models.Friendlist

	queryFriendlists := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", friendlistTable)

	if err := tx.Select(&friendlists, queryFriendlists, userID); err != nil {
		return nil, err
	}

	queryFriends := fmt.Sprintf(`SELECT f.* 
							FROM %s f
							INNER JOIN %s ff ON ff.friend_id = f.id 
							WHERE ff.friendlist_id = $1`, friendTable, friendlistsFriendsTable)

	stmt, err := tx.Preparex(queryFriends)
	if err != nil {
		return nil, err
	}

	var friendlistsWithFriends []models.FriendlistWithFriends

	var friends []models.Friend
	for _, friendlist := range friendlists {

		if err := stmt.Select(&friends, friendlist.ID); err != nil {
			return nil, err
		}
		friendlistsWithFriends = append(friendlistsWithFriends, models.FriendlistWithFriends{
			Friendlist: friendlist,
			Friends:    friends,
		})
		friends = nil
	}

	return friendlistsWithFriends, err

}

func (r *FriendlistPostgres) GetByIDWithFriends(userID, friendlistID uuid.UUID) (models.FriendlistWithFriends, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendlistWithFriends{}, err
	}
	defer tx.Rollback()

	var friendlist models.Friendlist

	queryFriendlist := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", friendlistTable)

	if err := tx.Get(&friendlist, queryFriendlist, friendlistID, userID); err != nil {
		return models.FriendlistWithFriends{}, err
	}

	queryFriends := fmt.Sprintf(`SELECT f.* 
							FROM %s f
							INNER JOIN %s ff ON ff.friend_id = f.id 
							WHERE ff.friendlist_id = $1`, friendTable, friendlistsFriendsTable)

	var friends []models.Friend
	if err := tx.Select(&friends, queryFriends, friendlist.ID); err != nil {
		return models.FriendlistWithFriends{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.FriendlistWithFriends{}, err
	}

	friendlistWithFriends := models.FriendlistWithFriends{
		Friendlist: friendlist,
		Friends:    friends,
	}

	return friendlistWithFriends, err
}

func (r *FriendlistPostgres) Update(userID, friendlistID uuid.UUID, friendlist models.UpdateFriendlist) error {

	friendlistFieldsWithValues := []string{}
	builderFriendlist := sqlbuilder.NewUpdateBuilder()
	builderFriendlist.SetFlavor(sqlbuilder.PostgreSQL)
	builderFriendlist.Update(friendlistTable)
	builderFriendlist.Where(
		builderFriendlist.Equal("id", friendlistID),
		builderFriendlist.Equal("user_id", userID),
	)

	if friendlist.Title != nil {
		friendlistFieldsWithValues = append(friendlistFieldsWithValues, builderFriendlist.Assign("title", *friendlist.Title))
	}
	if friendlist.Description != nil {
		friendlistFieldsWithValues = append(friendlistFieldsWithValues, builderFriendlist.Assign("description", *friendlist.Description))
	}
	if friendlist.Color != nil {
		friendlistFieldsWithValues = append(friendlistFieldsWithValues, builderFriendlist.Assign("color", *friendlist.Color))
	}
	if friendlist.ImageID != nil {
		friendlistFieldsWithValues = append(friendlistFieldsWithValues, builderFriendlist.Assign("image_id", *friendlist.ImageID))
	}

	builderFriendlist.Set(friendlistFieldsWithValues...)

	queryFriendlist, args := builderFriendlist.Build()

	_, err := r.db.Exec(queryFriendlist, args...)

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

func (r *FriendlistPostgres) AddTagsToFriendlist(userID, friendlistID uuid.UUID, tagIDs []models.AdditionTag) ([]uuid.UUID, error) {

	if len(tagIDs) == 0 {
		return nil, errors.New("empty tagIDs")
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO \"%s\" (friendlist_id, tag_id) VALUES ($1, $2) RETURNING id", friendlistsTagsTable)

	stmt, err := tx.Preparex(query)
	if err != nil {
		return nil, err
	}
	ids := []uuid.UUID{}
	var id uuid.UUID
	for _, tagID := range tagIDs {
		row := stmt.QueryRow(friendlistID, tagID.TagID)
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) != len(tagIDs) {
		return nil, errors.New("length of added tags not equal inserted rows")
	}

	err = tx.Commit()
	return ids, err
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
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friendlist_id = $1 AND tag_id = $2)", friendlistsTagsTable)

	if err := r.db.Get(&exists, queryCheck, friendlistID, tagID); err != nil {
		return err
	}
	if !exists {
		return errors.New("tag already removed from friendlist")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE friendlist_id=$1 AND tag_id=$2", friendlistsTagsTable)

	_, err := r.db.Exec(query, friendlistID, tagID)
	if err != nil {
		return err
	}

	return err
}

func (r *FriendlistPostgres) DeleteFriendFromFriendlist(friendlistID, friendID uuid.UUID) error {
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friendlist_id = $1 AND friend_id = $2)", friendlistsFriendsTable)

	if err := r.db.Get(&exists, queryCheck, friendlistID, friendID); err != nil {
		return err
	}
	if !exists {
		return errors.New("friend already removed from friendlist")
	}

	queryDelete := fmt.Sprintf("DELETE FROM %s where friendlist_id=$1 AND friend_id=$2", friendlistsFriendsTable)

	_, err := r.db.Exec(queryDelete, friendlistID, friendID)
	if err != nil {
		return err
	}

	return err
}
