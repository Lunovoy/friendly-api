package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type ReminderService struct {
	repo repository.Reminder
}

func NewReminderService(repo repository.Reminder) *ReminderService {
	return &ReminderService{
		repo: repo,
	}
}

func (s *ReminderService) Create(userID uuid.UUID, reminder models.Reminder) (uuid.UUID, error) {
	return s.repo.Create(userID, reminder)
}

func (s *ReminderService) CreateBulk(userID, eventID uuid.UUID, reminders []models.Reminder) ([]uuid.UUID, error) {
	return s.repo.CreateBulk(userID, eventID, reminders)
}

func (s *ReminderService) GetAll(userID uuid.UUID) ([]models.Reminder, error) {
	return s.repo.GetAll(userID)
}

func (s *ReminderService) GetAllByEventID(userID, eventID uuid.UUID) ([]models.Reminder, error) {
	return s.repo.GetAllByEventID(userID, eventID)
}

func (s *ReminderService) GetByID(userID, reminderID uuid.UUID) (models.Reminder, error) {
	return s.repo.GetByID(userID, reminderID)
}

func (s *ReminderService) DeleteByID(userID, reminderID uuid.UUID) error {
	return s.repo.DeleteByID(userID, reminderID)
}
