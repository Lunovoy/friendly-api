package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type EventService struct {
	repo repository.Event
}

func NewEventService(repo repository.Event) *EventService {
	return &EventService{
		repo: repo,
	}
}

var frequencies = map[string]bool{
	"once":        true,
	"everyday":    true,
	"weekdays":    true, //будние
	"weekly":      true, //еженедельно (через 7 дней)
	"monthlyDate": true, // ежемесячно (в эту же дату)
	"monthlyDay":  true, //ежемесячно (в 4й вторник к примеру)
	"annually":    true, //ежегодно
}

func (s *EventService) isFrequencyValid(frequency string) bool {
	if _, ok := frequencies[frequency]; !ok {
		return false
	}
	return true
}

func (s *EventService) Create(userID uuid.UUID, event models.Event) (uuid.UUID, error) {
	if !s.isFrequencyValid(event.Frequency) {
		return uuid.Nil, errors.New("frequency is not valid")
	}
	if !event.EndDate.Valid {
		event.EndDate.Time = event.StartDate.Time.Add(5 * time.Minute)
		event.EndDate.Valid = true
	}
	return s.repo.Create(userID, event)
}

func (s *EventService) AddFriendsToEvent(userID, eventID uuid.UUID, friendIDs []models.FriendID) ([]uuid.UUID, error) {
	return s.repo.AddFriendsToEvent(userID, eventID, friendIDs)
}

func (s *EventService) DeleteFriendsFromEvent(userID, eventID uuid.UUID, friendIDs []uuid.UUID) error {
	return s.repo.DeleteFriendsFromEvent(userID, eventID, friendIDs)
}

func (s *EventService) GetEventsByFriendID(userID, friendID uuid.UUID) ([]models.Event, error) {
	return s.repo.GetEventsByFriendID(userID, friendID)
}

func (s *EventService) GetAll(userID uuid.UUID) ([]models.Event, error) {
	return s.repo.GetAll(userID)
}

func (s *EventService) GetByID(userID, eventID uuid.UUID) (models.Event, error) {
	return s.repo.GetByID(userID, eventID)
}

func (s *EventService) GetAllWithFriends(userID uuid.UUID) ([]models.EventWithFriends, error) {
	return s.repo.GetAllWithFriends(userID)
}

func (s *EventService) GetByIDWithFriends(userID, eventID uuid.UUID) (models.EventWithFriends, error) {
	return s.repo.GetByIDWithFriends(userID, eventID)
}

func (s *EventService) Update(userID, eventID uuid.UUID, event models.EventUpdate) error {
	if event.Frequency != nil {
		if !s.isFrequencyValid(*event.Frequency) {
			return errors.New("frequency is not valid")
		}
	}
	return s.repo.Update(userID, eventID, event)
}

func (s *EventService) UpdateFull(userID, eventID uuid.UUID, event models.EventFullUpdate) error {
	if event.EventUpdate != nil {
		if event.EventUpdate.Frequency != nil {
			if !s.isFrequencyValid(*event.EventUpdate.Frequency) {
				return errors.New("frequency is not valid")
			}
		}
	}
	return s.repo.UpdateFull(userID, eventID, event)
}

func (s *EventService) DeleteByID(userID, eventID uuid.UUID) error {
	return s.repo.DeleteByID(userID, eventID)
}
