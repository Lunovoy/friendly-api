package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type FriendService struct {
	repo repository.Friend
}

func NewFriendService(repo repository.Friend) *FriendService {
	return &FriendService{
		repo: repo,
	}
}

func (s *FriendService) Create(userID uuid.UUID, friend models.Friend, workInfo models.WorkInfo) (models.FriendIDWorkInfoID, error) {
	return s.repo.Create(userID, friend, workInfo)
}

func (s *FriendService) GetAll(userID uuid.UUID) ([]models.FriendWorkInfo, error) {
	return s.repo.GetAll(userID)
}

func (s *FriendService) GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfo, error) {
	return s.repo.GetByID(userID, friendID)
}

func (s *FriendService) Update(userID, friendID uuid.UUID, friend models.Friend) error {
	return s.repo.Update(userID, friendID, friend)
}

func (s *FriendService) DeleteByID(userID, friendID uuid.UUID) error {
	return s.repo.DeleteByID(userID, friendID)
}
