package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type ReminderPostgres struct {
	db *sqlx.DB
}

func NewReminderPostgres(db *sqlx.DB) *ReminderPostgres {
	return &ReminderPostgres{
		db: db,
	}
}

func (r *ReminderPostgres) Create(userID uuid.UUID, reminder models.Reminder) (uuid.UUID, error) {
	var reminderID uuid.UUID
	query := fmt.Sprintf("INSERT INTO \"%s\" (minutes_until_event, event_id, user_id) VALUES ($1, $2, $3) RETURNING id", reminderTable)

	row := r.db.QueryRow(query, reminder.MinutesUntilEvent, reminder.EventID, userID)
	if err := row.Scan(&reminderID); err != nil {
		return uuid.Nil, err
	}

	return reminderID, nil
}

func (r *ReminderPostgres) CreateBulk(userID, eventID uuid.UUID, reminders []models.Reminder) ([]uuid.UUID, error) {

	query := fmt.Sprintf("INSERT INTO \"%s\" (minutes_until_event, event_id, user_id) VALUES ($1, $2, $3) RETURNING id", reminderTable)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var reminderID *uuid.UUID
	var reminderIDs []uuid.UUID
	for _, reminder := range reminders {
		row := stmt.QueryRow(reminder.MinutesUntilEvent, eventID, userID)
		if err := row.Scan(&reminderID); err != nil {
			return nil, err
		}
		reminderIDs = append(reminderIDs, *reminderID)
		reminderID = nil

	}

	return reminderIDs, nil
}

func (r *ReminderPostgres) GetAll(userID uuid.UUID) ([]models.Reminder, error) {

	var reminders []models.Reminder

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", reminderTable)

	err := r.db.Select(&reminders, query, userID)

	return reminders, err
}

func (r *ReminderPostgres) GetAllByEventID(userID, eventID uuid.UUID) ([]models.Reminder, error) {

	var reminders []models.Reminder

	query := fmt.Sprintf("SELECT * FROM %s WHERE event_id = $1, user_id = $2", reminderTable)

	err := r.db.Select(&reminders, query, eventID, userID)

	return reminders, err
}

func (r *ReminderPostgres) GetByID(userID, reminderID uuid.UUID) (models.Reminder, error) {
	var reminder models.Reminder

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", reminderTable)

	err := r.db.Get(&reminder, query, reminderID, userID)

	return reminder, err
}

func (r *ReminderPostgres) DeleteByID(userID, reminderID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND user_id=$2", reminderID)

	_, err := r.db.Exec(query, reminderID, userID)

	return err
}
