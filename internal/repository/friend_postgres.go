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

func (r *FriendPostgres) Create(userID uuid.UUID, friend models.Friend, workInfo models.WorkInfo) (models.FriendIDWorkInfoID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return models.FriendIDWorkInfoID{}, err
	}
	defer tx.Rollback()

	var friendID uuid.UUID
	queryFriend := fmt.Sprintf("INSERT INTO \"%s\" (first_name, last_name , dob, user_id) VALUES ($1, $2, $3, $4) RETURNING id", friendTable)

	rowFriend := tx.QueryRow(queryFriend, friend.FirstName, friend.LastName, friend.DOB, userID)
	if err := rowFriend.Scan(&friendID); err != nil {
		return models.FriendIDWorkInfoID{}, err
	}

	var workInfoID uuid.UUID
	queryWorkInfo := fmt.Sprintf("INSERT INTO \"%s\" (country, city, company, position, messenger, communication_method, nationality, language, friend_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", workInfoTable)

	rowWorkInfo := tx.QueryRow(queryWorkInfo, workInfo.Country, workInfo.City, workInfo.Company, workInfo.Position, workInfo.Messenger, workInfo.CommunicationMethod, workInfo.Nationality, workInfo.Language, friendID)
	if err := rowWorkInfo.Scan(&workInfoID); err != nil {
		return models.FriendIDWorkInfoID{}, err
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
		builderWorkInfo.Update(friendTable)
		builderWorkInfo.Where(
			builderWorkInfo.Equal("id", friendID),
			builderWorkInfo.Equal("user_id", userID),
		)

		if friend.Friend.FirstName != nil {
			builderWorkInfo.Set("first_name", *friend.Friend.FirstName)
		}
		if friend.Friend.LastName != nil {
			builderWorkInfo.Set("last_name", *friend.Friend.LastName)
		}
		if friend.Friend.DOB != nil {
			builderWorkInfo.Set("dob", friend.Friend.DOB.Format("2006-01-02 15:04:05-07:00"))
		}

		queryFriend, args := builderWorkInfo.Build()
		_, err = r.db.Exec(queryFriend, args)
		if err != nil {
			return err
		}
	}

	queryWorkInfo := fmt.Sprintf("UPDATE %s SET country = $1, city = $2, company = $3, position = $4, messenger = $5, communication_method = $6, nationality = $7, language WHERE friend_id = $8", workInfoTable)

	_, err = r.db.Exec(queryWorkInfo, friend.WorkInfo.Country, friend.WorkInfo.City, friend.WorkInfo.Company, friend.WorkInfo.Position, friend.WorkInfo.Messenger, friend.WorkInfo.CommunicationMethod, friend.WorkInfo.Nationality, friendID)
	if err != nil {
		return err
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
