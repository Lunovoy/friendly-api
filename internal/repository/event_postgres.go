package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type EventPostgres struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgres {
	return &EventPostgres{
		db: db,
	}
}

func (r *EventPostgres) Create(userID uuid.UUID, event models.Event) (uuid.UUID, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var eventID uuid.UUID
	query := fmt.Sprintf("INSERT INTO \"%s\" (title, description, start_date, end_date, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", eventTable)

	row := tx.QueryRow(query, event.Title, event.Description, event.StartDate, event.EndDate, userID)
	if err := row.Scan(&eventID); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return eventID, nil
}

func (r *EventPostgres) AddFriendsToEvent(userID, eventID uuid.UUID, friendIDs []uuid.UUID) ([]uuid.UUID, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO \"%s\" (friend_id, event_id) VALUES ($1, $2) RETURNING id", friendsEventsTable)

	stmt, err := tx.Preparex(query)
	if err != nil {
		return nil, err
	}
	ids := []uuid.UUID{}
	var id uuid.UUID
	for _, friendID := range friendIDs {
		row := stmt.QueryRow(stmt, friendID, eventID)
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) != len(friendIDs) {
		return nil, errors.New("length of added friends not equal inserted rows")
	}

	err = tx.Commit()
	return ids, err
}

func (r *EventPostgres) GetEventsByFriendID(userID, friendID uuid.UUID) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf(`SELECT e.*
						FROM %s e
						JOIN %s fe ON fe.event_id = e.id
						WHERE fe.friend_id = $1 AND e.user_id = $2`, eventTable, friendsEventsTable)

	err := r.db.Select(&events, query, friendID, userID)

	return events, err
}

func (r *EventPostgres) GetAll(userID uuid.UUID) ([]models.Event, error) {
	var events []models.Event

	query := fmt.Sprintf("SELECT * FROM %s where user_id = $1", eventTable)

	err := r.db.Select(&events, query, userID)

	return events, err

}

func (r *EventPostgres) GetByID(userID, eventID uuid.UUID) (models.Event, error) {
	var event models.Event

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", eventTable)

	err := r.db.Get(&event, query, eventID, userID)

	return event, err
}

func (r *EventPostgres) Update(userID, eventID uuid.UUID, event models.Event) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5 AND user_id = $6", eventTable)

	_, err := r.db.Exec(query, event.Title, event.Description, event.StartDate, event.EndDate, eventID, userID)

	return err
}

func (r *EventPostgres) DeleteByID(userID, eventID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND user_id=$2", eventTable)

	_, err := r.db.Exec(query, eventID, userID)

	return err
}
