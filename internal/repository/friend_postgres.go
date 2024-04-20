package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type FriendPostgres struct {
	db *sqlx.DB
}

func NewFriendPostgres(db *sqlx.DB) *FriendPostgres {
	return &FriendPostgres{
		db: db,
	}
}

func (r *FriendPostgres) Create(userID uuid.UUID, friend models.UpdateFriendWorkInfoInput) (models.FriendIDWorkInfoID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendIDWorkInfoID{}, err
	}
	defer tx.Rollback()

	friendFields := []string{"first_name", "user_id"}
	friendValues := []any{*friend.Friend.FirstName, userID}
	builderFriend := sqlbuilder.NewInsertBuilder()
	builderFriend.SetFlavor(sqlbuilder.PostgreSQL)
	builderFriend.InsertInto(friendTable)

	if friend.Friend.LastName != nil {
		friendFields = append(friendFields, "last_name")
		friendValues = append(friendValues, *friend.Friend.LastName)
	}
	if friend.Friend.DOB != nil {
		friendFields = append(friendFields, "dob")
		friendValues = append(friendValues, *friend.Friend.DOB)
	}
	if friend.Friend.ImageID != nil {
		friendFields = append(friendFields, "image_id")
		friendValues = append(friendValues, *friend.Friend.ImageID)
	}

	builderFriend.Cols(friendFields...).Values(friendValues...)

	queryFriend, args := builderFriend.Build()
	queryFriend += " RETURNING id;"

	var friendID uuid.UUID

	rowFriend := tx.QueryRow(queryFriend, args...)
	if err := rowFriend.Scan(&friendID); err != nil {
		return models.FriendIDWorkInfoID{}, err
	}

	var workInfoID uuid.UUID
	builderWorkInfo := sqlbuilder.NewInsertBuilder()
	builderWorkInfo.InsertInto(workInfoTable)
	builderWorkInfo.SetFlavor(sqlbuilder.PostgreSQL)

	if friend.WorkInfo != nil {

		workFields := []string{"friend_id"}
		workValues := []any{friendID}

		fieldsToUpdateWorkInfo := map[string]*string{
			"country":              friend.WorkInfo.Country,
			"city":                 friend.WorkInfo.City,
			"company":              friend.WorkInfo.Company,
			"position":             friend.WorkInfo.Position,
			"messenger":            friend.WorkInfo.Messenger,
			"communication_method": friend.WorkInfo.CommunicationMethod,
			"nationality":          friend.WorkInfo.Nationality,
			"language":             friend.WorkInfo.Language,
		}

		for field, value := range fieldsToUpdateWorkInfo {
			if value != nil {
				workFields = append(workFields, field)
				workValues = append(workValues, value)
			}
		}
		builderWorkInfo.Cols(workFields...).Values(workValues...)

		queryWorkInfo, args := builderWorkInfo.Build()
		queryWorkInfo += " RETURNING id;"

		rowWorkInfo := tx.QueryRow(queryWorkInfo, args...)

		if err := rowWorkInfo.Scan(&workInfoID); err != nil {
			return models.FriendIDWorkInfoID{}, err
		}
	} else {

		builderWorkInfo.Cols("friend_id").Values(friendID)

		queryWorkInfo, args := builderWorkInfo.Build()
		queryWorkInfo += " RETURNING id;"
		rowWorkInfo := tx.QueryRow(queryWorkInfo, args...)

		if err := rowWorkInfo.Scan(&workInfoID); err != nil {
			return models.FriendIDWorkInfoID{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.FriendIDWorkInfoID{}, err
	}

	return models.FriendIDWorkInfoID{FriendID: friendID, WorkInfoID: workInfoID}, nil
}

func (r *FriendPostgres) GetAll(userID uuid.UUID) ([]models.FriendWorkInfoTags, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var friends []models.Friend
	friendQuery := fmt.Sprintf("SELECT id, first_name, last_name, dob, image_id, user_id FROM %s WHERE user_id = $1", friendTable)
	err = tx.Select(&friends, friendQuery, userID)
	if err != nil {
		return nil, err
	}

	var workInfos []models.WorkInfo
	workInfoQuery := fmt.Sprintf(`SELECT w.id, w.country, w.city, w.company, w.position, 
								w.messenger, w.communication_method, w.nationality, 
								w.language, w.friend_id 
								FROM %s w
								INNER JOIN %s f ON w.friend_id = f.id
    							WHERE f.user_id = $1`, workInfoTable, friendTable)
	err = tx.Select(&workInfos, workInfoQuery, userID)
	if err != nil {
		return nil, err
	}

	queryTags := fmt.Sprintf(`SELECT t.id, t.title, t.user_id 
							FROM %s t
							INNER JOIN %s ft ON ft.tag_id = t.id 
							WHERE ft.friend_id = $1`, tagTable, friendsTagsTable)

	tagStmt, err := tx.Preparex(queryTags)
	if err != nil {
		return nil, err
	}

	var tags []models.Tag

	var friendWorkInfos []models.FriendWorkInfoTags
	for _, friend := range friends {
		if err := tagStmt.Select(&tags, friend.ID); err != nil {
			return nil, err
		}
		for _, workInfo := range workInfos {
			if friend.ID == workInfo.FriendID {
				friendWorkInfos = append(friendWorkInfos, models.FriendWorkInfoTags{
					Friend:   friend,
					WorkInfo: workInfo,
					Tags:     tags,
				})
			}
		}
		tags = nil
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return friendWorkInfos, nil
}

func (r *FriendPostgres) GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfoTags, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendWorkInfoTags{}, err
	}
	defer tx.Rollback()

	var friend models.Friend
	friendQuery := fmt.Sprintf("SELECT id, first_name, last_name, dob, image_id, user_id FROM %s WHERE id = $1 AND user_id = $2", friendTable)
	err = tx.Get(&friend, friendQuery, friendID, userID)
	if err != nil {
		return models.FriendWorkInfoTags{}, err
	}

	var workInfo models.WorkInfo
	workInfoQuery := fmt.Sprintf(`SELECT w.id, w.country, w.city, w.company, w.position, 
								w.messenger, w.communication_method, w.nationality, 
								w.language , w.friend_id
								FROM %s w
								INNER JOIN %s f ON w.friend_id = $1 
    							WHERE f.user_id = $2`, workInfoTable, friendTable)
	err = tx.Get(&workInfo, workInfoQuery, friendID, userID)
	if err != nil {
		return models.FriendWorkInfoTags{}, err
	}

	queryTags := fmt.Sprintf(`SELECT t.id, t.title, t.user_id 
	FROM %s t
	INNER JOIN %s ft ON ft.tag_id = t.id 
	WHERE ft.friend_id = $1`, tagTable, friendsTagsTable)

	var tags []models.Tag
	if err := tx.Select(&tags, queryTags, friend.ID); err != nil {
		return models.FriendWorkInfoTags{}, err
	}

	friendWorkInfo := models.FriendWorkInfoTags{
		Friend:   friend,
		WorkInfo: workInfo,
		Tags:     tags,
	}

	if err := tx.Commit(); err != nil {
		return models.FriendWorkInfoTags{}, err
	}

	return friendWorkInfo, nil
}

func (r *FriendPostgres) Update(userID, friendID uuid.UUID, friend models.UpdateFriendWorkInfoInput) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if friend.Friend != nil {
		friendFieldsWithValues := []string{}
		builderFriend := sqlbuilder.NewUpdateBuilder()
		builderFriend.SetFlavor(sqlbuilder.PostgreSQL)
		builderFriend.Update(friendTable)
		builderFriend.Where(
			builderFriend.Equal("id", friendID),
			builderFriend.Equal("user_id", userID),
		)

		if friend.Friend.FirstName != nil {
			friendFieldsWithValues = append(friendFieldsWithValues, builderFriend.Assign("first_name", *friend.Friend.FirstName))
		}
		if friend.Friend.LastName != nil {
			friendFieldsWithValues = append(friendFieldsWithValues, builderFriend.Assign("last_name", *friend.Friend.LastName))
		}
		if friend.Friend.DOB != nil {
			friendFieldsWithValues = append(friendFieldsWithValues, builderFriend.Assign("dob", *friend.Friend.DOB))
		}
		if friend.Friend.ImageID != nil {
			friendFieldsWithValues = append(friendFieldsWithValues, builderFriend.Assign("image_id", *friend.Friend.ImageID))
		}

		builderFriend.Set(friendFieldsWithValues...)

		queryFriend, args := builderFriend.Build()
		_, err = r.db.Exec(queryFriend, args...)
		if err != nil {
			return err
		}
	}

	if friend.WorkInfo != nil {
		workFieldsWithValues := []string{}
		builderWorkInfo := sqlbuilder.NewUpdateBuilder()
		builderWorkInfo.SetFlavor(sqlbuilder.PostgreSQL)
		builderWorkInfo.Update(workInfoTable)
		builderWorkInfo.Where(
			builderWorkInfo.Equal("friend_id", friendID),
		)

		fieldsToUpdateWorkInfo := map[string]*string{
			"country":              friend.WorkInfo.Country,
			"city":                 friend.WorkInfo.City,
			"company":              friend.WorkInfo.Company,
			"position":             friend.WorkInfo.Position,
			"messenger":            friend.WorkInfo.Messenger,
			"communication_method": friend.WorkInfo.CommunicationMethod,
			"nationality":          friend.WorkInfo.Nationality,
			"language":             friend.WorkInfo.Language,
		}

		for field, value := range fieldsToUpdateWorkInfo {
			if value != nil {
				workFieldsWithValues = append(workFieldsWithValues, builderWorkInfo.Assign(field, *value))
			}
		}
		builderWorkInfo.Set(workFieldsWithValues...)

		queryWorkInfo, args := builderWorkInfo.Build()
		_, err = r.db.Exec(queryWorkInfo, args...)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (r *FriendPostgres) DeleteByID(userID, FriendID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s where id=$1 AND user_id=$2", friendTable)

	_, err := r.db.Exec(query, FriendID, userID)
	if err != nil {
		return err
	}

	return err
}

func (r *FriendPostgres) AddTagToFriend(friendID, tagID uuid.UUID) error {

	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friend_id = $1 AND tag_id = $2)", friendsTagsTable)

	if err := r.db.Get(&exists, queryCheck, friendID, tagID); err != nil {
		return err
	}
	if exists {
		return errors.New("tag already exists for this friend")
	}

	queryAdd := fmt.Sprintf("INSERT INTO \"%s\" (friend_id, tag_id) VALUES ($1, $2)", friendsTagsTable)

	_, err := r.db.Exec(queryAdd, friendID, tagID)

	return err
}

func (r *FriendPostgres) DeleteTagFromFriend(friendID, tagID uuid.UUID) error {
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE friend_id = $1 AND tag_id = $2)", friendsTagsTable)

	if err := r.db.Get(&exists, queryCheck, friendID, tagID); err != nil {
		return err
	}
	if !exists {
		return errors.New("tag already removed from friend")
	}

	query := fmt.Sprintf("DELETE FROM %s where friend_id=$1 AND tag_id=$2", friendsTagsTable)

	_, err := r.db.Exec(query, friendID, tagID)
	if err != nil {
		return err
	}

	return err
}
