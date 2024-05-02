package service

import (
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

func (s *EventService) Create(userID uuid.UUID, event models.Event) (uuid.UUID, error) {
	return s.repo.Create(userID, event)
}

func (s *EventService) AddFriendsToEvent(userID, eventID uuid.UUID, friendIDs []models.FriendID) ([]uuid.UUID, error) {
	return s.repo.AddFriendsToEvent(userID, eventID, friendIDs)
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
	return s.repo.Update(userID, eventID, event)
}

func (s *EventService) DeleteByID(userID, eventID uuid.UUID) error {
	return s.repo.DeleteByID(userID, eventID)
}
