package repository

import (
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

	fmt.Println(*friend.Friend.FirstName)

	builderFriend := sqlbuilder.NewInsertBuilder()
	builderFriend.InsertInto(friendTable)

	if friend.Friend.LastName != nil {
		builderFriend.Cols("last_name").Values(*friend.Friend.LastName)
	}
	if friend.Friend.DOB != nil {
		builderFriend.Cols("dob").Values(*friend.Friend.DOB)
	}

	builderFriend.Cols("first_name").Values(*friend.Friend.FirstName)
	// builderFriend.Cols("user_id").Values(userID)

	// queryFriend := builderFriend.String() + " RETURNING id;"
	queryFriend, args := builderFriend.Build()
	fmt.Println(queryFriend, args)

	var friendID uuid.UUID

	rowFriend := tx.QueryRow(queryFriend)
	if err := rowFriend.Scan(&friendID); err != nil {
		return models.FriendIDWorkInfoID{}, err
	}

	var workInfoID uuid.UUID
	builderWorkInfo := sqlbuilder.NewInsertBuilder()
	builderWorkInfo.InsertInto(workInfoTable)

	if friend.WorkInfo != nil {

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
				builderWorkInfo.Cols(field).Values(*value)
			}
		}

		builderWorkInfo.Cols("friend_id").Values(friendID)

		queryWorkInfo := builderWorkInfo.String() + " RETURNING id;"

		rowWorkInfo := tx.QueryRow(queryWorkInfo)

		if err := rowWorkInfo.Scan(&workInfoID); err != nil {
			return models.FriendIDWorkInfoID{}, err
		}
	} else {

		builderWorkInfo.Cols("friend_id").Values(friendID)

		queryWorkInfo := builderWorkInfo.String() + " RETURNING id;"
		rowWorkInfo := tx.QueryRow(queryWorkInfo)

		if err := rowWorkInfo.Scan(&workInfoID); err != nil {
			return models.FriendIDWorkInfoID{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.FriendIDWorkInfoID{}, err
	}

	return models.FriendIDWorkInfoID{FriendID: friendID, WorkInfoID: workInfoID}, nil
}

func (r *FriendPostgres) GetAll(userID uuid.UUID) ([]models.FriendWorkInfo, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var friends []models.Friend
	friendQuery := fmt.Sprintf("SELECT id, first_name, last_name, dob, user_id FROM %s WHERE user_id = $1", friendTable)
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

	var friendWorkInfos []models.FriendWorkInfo
	for _, friend := range friends {
		for _, workInfo := range workInfos {
			if friend.ID == workInfo.FriendID {
				friendWorkInfos = append(friendWorkInfos, models.FriendWorkInfo{
					Friend:   friend,
					WorkInfo: workInfo,
				})
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return friendWorkInfos, nil
}

func (r *FriendPostgres) GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfo, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendWorkInfo{}, err
	}
	defer tx.Rollback()

	var friend models.Friend
	friendQuery := fmt.Sprintf("SELECT id, first_name, last_name, dob, user_id FROM %s WHERE id = $1 AND user_id = $2", friendTable)
	err = tx.Get(&friend, friendQuery, friendID, userID)
	if err != nil {
		return models.FriendWorkInfo{}, err
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
		return models.FriendWorkInfo{}, err
	}

	friendWorkInfo := models.FriendWorkInfo{
		Friend:   friend,
		WorkInfo: workInfo,
	}

	if err := tx.Commit(); err != nil {
		return models.FriendWorkInfo{}, err
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
		builderFriend := sqlbuilder.NewUpdateBuilder()
		builderFriend.Update(friendTable)
		builderFriend.Where(
			builderFriend.Equal("id", friendID),
			builderFriend.Equal("user_id", userID),
		)

		if friend.Friend.FirstName != nil {
			builderFriend.Set("first_name", *friend.Friend.FirstName)
		}
		if friend.Friend.LastName != nil {
			builderFriend.Set("last_name", *friend.Friend.LastName)
		}
		if friend.Friend.DOB != nil {
			builderFriend.Set("dob", friend.Friend.DOB.Format("2006-01-02 15:04:05-07:00"))
		}

		queryFriend, args := builderFriend.Build()
		_, err = r.db.Exec(queryFriend, args)
		if err != nil {
			return err
		}
	}

	if friend.WorkInfo != nil {
		builderWorkInfo := sqlbuilder.NewUpdateBuilder()
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
				builderWorkInfo.Set(field, *value)
			}
		}

		queryWorkInfo, args := builderWorkInfo.Build()
		_, err = r.db.Exec(queryWorkInfo, args)
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
