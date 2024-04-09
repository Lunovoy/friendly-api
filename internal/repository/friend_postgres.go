package repository

import (
	"fmt"

	"github.com/google/uuid"
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
	queryWorkInfo := fmt.Sprintf("INSERT INTO \"%s\" (country, city, company, position, messenger, communication_method, nationality, friend_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", workInfoTable)

	rowWorkInfo := tx.QueryRow(queryWorkInfo, workInfo.Country, workInfo.City, workInfo.Company, workInfo.Position, workInfo.Messenger, workInfo.CommunicationMethod, workInfo.Nationality, friendID)
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
								w.messenger, w.communication_method, w.nationality, w.friend_id 
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
								w.messenger, w.communication_method, w.nationality, w.friend_id 
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

func (r *FriendPostgres) Update(userID, FriendID uuid.UUID, friend models.Friend) error {
	query := fmt.Sprintf("UPDATE %s SET first_name = $1, last_name = $2, dob = $3 WHERE id = $4 AND user_id = $5", friendTable)

	_, err := r.db.Exec(query, friend.FirstName, friend.LastName, friend.DOB, FriendID, userID)

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
