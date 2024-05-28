package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
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
	query := fmt.Sprintf("INSERT INTO \"%s\" (title, description, start_date, end_date, frequency, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", eventTable)

	row := tx.QueryRow(query, event.Title, event.Description, event.StartDate, event.EndDate, event.Frequency, userID)
	if err := row.Scan(&eventID); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return eventID, nil
}

func (r *EventPostgres) AddFriendsToEvent(userID, eventID uuid.UUID, friendIDs []models.FriendID) ([]uuid.UUID, error) {

	if len(friendIDs) == 0 {
		return nil, nil
	}

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
		row := stmt.QueryRow(friendID.FriendID, eventID)
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

func (r *EventPostgres) DeleteFriendsFromEvent(userID, eventID uuid.UUID, friendIDs []uuid.UUID) error {
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE event_id = $1 AND friend_id = $2)", friendsEventsTable)

	checkStmt, err := r.db.Preparex(queryCheck)
	if err != nil {
		return err
	}
	defer checkStmt.Close()

	for _, friendID := range friendIDs {
		if err := checkStmt.Get(&exists, eventID, friendID); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("friend %s already removed from event", friendID)
		}
	}

	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE event_id=$1 AND friend_id=$2", friendsEventsTable)

	deleteStmt, err := r.db.Preparex(queryDelete)
	if err != nil {
		return err
	}

	for _, friendID := range friendIDs {
		_, err := deleteStmt.Exec(eventID, friendID)
		if err != nil {
			return err
		}
	}

	return err
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

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", eventTable)

	err := r.db.Select(&events, query, userID)

	return events, err

}

func (r *EventPostgres) GetByID(userID, eventID uuid.UUID) (models.Event, error) {
	var event models.Event

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", eventTable)

	fmt.Println("Event: ", eventID, "User: ", userID)
	err := r.db.Get(&event, query, eventID, userID)

	return event, err
}

func (r *EventPostgres) GetByIDWithFriends(userID, eventID uuid.UUID) (models.EventWithFriends, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return models.EventWithFriends{}, err
	}
	defer tx.Rollback()

	queryEvent := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", eventTable)

	var event models.Event

	err = tx.Get(&event, queryEvent, eventID, userID)
	if err != nil {
		return models.EventWithFriends{}, err
	}

	queryFriends := fmt.Sprintf(`SELECT f.* 
								FROM %s f
								JOIN %s fe ON f.id = fe.friend_id
								WHERE fe.event_id = $1 AND f.user_id = $2`, friendTable, friendsEventsTable)

	var friends []models.Friend

	err = tx.Select(&friends, queryFriends, eventID, userID)
	if err != nil {
		return models.EventWithFriends{}, err
	}

	eventWithFriends := models.EventWithFriends{
		Event:   event,
		Friends: friends,
	}

	err = tx.Commit()
	if err != nil {
		return models.EventWithFriends{}, err
	}

	return eventWithFriends, nil
}

func (r *EventPostgres) GetAllWithFriends(userID uuid.UUID) ([]models.EventWithFriends, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	queryEvents := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", eventTable)

	var events []models.Event
	err = tx.Select(&events, queryEvents, userID)
	if err != nil {
		return nil, err
	}

	var eventWithFriendsList []models.EventWithFriends

	queryFriends := fmt.Sprintf(`SELECT f.* 
								FROM %s f
								JOIN %s fe ON f.id = fe.friend_id
								WHERE fe.event_id = $1 AND f.user_id = $2`, friendTable, friendsEventsTable)

	friendsStmt, err := tx.Preparex(queryFriends)
	if err != nil {
		return nil, err
	}
	defer friendsStmt.Close()

	for _, event := range events {

		var friends []models.Friend
		err = friendsStmt.Select(&friends, event.ID, userID)
		if err != nil {
			return nil, err
		}

		eventWithFriends := models.EventWithFriends{
			Event:   event,
			Friends: friends,
		}

		eventWithFriendsList = append(eventWithFriendsList, eventWithFriends)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return eventWithFriendsList, nil

}

func (r *EventPostgres) Update(userID, eventID uuid.UUID, event models.EventUpdate) error {

	eventFieldsWithValues := []string{}
	builderEvent := sqlbuilder.NewUpdateBuilder()
	builderEvent.SetFlavor(sqlbuilder.PostgreSQL)
	builderEvent.Update(eventTable)
	builderEvent.Where(
		builderEvent.Equal("id", eventID),
		builderEvent.Equal("user_id", userID),
	)

	if event.Title != nil {
		eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("title", *event.Title))
	}
	if event.Description != nil {
		eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("description", *event.Description))
	}
	if event.StartDate != nil {
		eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("start_date", *event.StartDate))
	}
	if event.EndDate != nil {
		eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("end_date", *event.EndDate))
	}
	if event.Frequency != nil {
		eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("frequency", *event.Frequency))
	}

	builderEvent.Set(eventFieldsWithValues...)

	queryEvent, args := builderEvent.Build()
	_, err := r.db.Exec(queryEvent, args...)
	if err != nil {
		return err
	}

	return err
}

func (r *EventPostgres) UpdateFull(userID, eventID uuid.UUID, event models.EventFullUpdate) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if event.EventUpdate != nil {

		eventFieldsWithValues := []string{}
		builderEvent := sqlbuilder.NewUpdateBuilder()
		builderEvent.SetFlavor(sqlbuilder.PostgreSQL)
		builderEvent.Update(eventTable)
		builderEvent.Where(
			builderEvent.Equal("id", eventID),
			builderEvent.Equal("user_id", userID),
		)

		if event.EventUpdate.Title != nil {
			eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("title", *event.EventUpdate.Title))
		}
		if event.EventUpdate.Description != nil {
			eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("description", *event.EventUpdate.Description))
		}
		if event.EventUpdate.StartDate != nil {
			eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("start_date", *event.EventUpdate.StartDate))
		}
		if event.EventUpdate.EndDate != nil {
			eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("end_date", *event.EventUpdate.EndDate))
		}
		if event.EventUpdate.Frequency != nil {
			eventFieldsWithValues = append(eventFieldsWithValues, builderEvent.Assign("frequency", *event.EventUpdate.Frequency))
		}

		builderEvent.Set(eventFieldsWithValues...)

		queryEvent, args := builderEvent.Build()
		_, err := tx.Exec(queryEvent, args...)
		if err != nil {
			return err
		}
	}

	if event.FriendIDs != nil {
		queryDelete := fmt.Sprintf("DELETE FROM %s WHERE event_id = $1", friendsEventsTable)
		_, err := tx.Exec(queryDelete, eventID)
		if err != nil {
			return fmt.Errorf("error deleting old friends from event: %s", err.Error())
		}

		var friendsIDs []models.FriendID
		for _, friendID := range event.FriendIDs {
			friendsIDs = append(friendsIDs, *friendID)
		}
		_, err = r.AddFriendsToEvent(userID, eventID, friendsIDs)
		if err != nil {
			return fmt.Errorf("error adding new friends to event: %s", err.Error())
		}
	}

	if event.ReminderUpdate != nil {
		queryReminderDelete := fmt.Sprintf("DELETE FROM %s WHERE event_id=$1", reminderTable)
		_, err := tx.Exec(queryReminderDelete, eventID)
		if err != nil {
			return fmt.Errorf("error deleting old reminders from event: %s", err.Error())
		}

		queryAddNew := fmt.Sprintf("INSERT INTO \"%s\" (minutes_until_event, event_id, user_id) VALUES ($1, $2, $3) RETURNING id", reminderTable)

		reminderStmt, err := tx.Preparex(queryAddNew)
		if err != nil {
			return err
		}
		defer reminderStmt.Close()

		var reminderID uuid.UUID
		for _, reminder := range event.ReminderUpdate {
			row := reminderStmt.QueryRow(reminder.MinutesUntilEvent, eventID, userID)
			if err := row.Scan(&reminderID); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (r *EventPostgres) DeleteByID(userID, eventID uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND user_id=$2", eventTable)

	_, err := r.db.Exec(query, eventID, userID)

	return err
}
