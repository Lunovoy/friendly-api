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

func (s *FriendService) Create(userID uuid.UUID, friend models.UpdateFriendWorkInfoInput) (models.FriendIDWorkInfoID, error) {
	return s.repo.Create(userID, friend)
}

func (s *FriendService) GetAll(userID uuid.UUID) ([]models.FriendWorkInfoTags, error) {
	return s.repo.GetAll(userID)
}

func (s *FriendService) GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfoTags, error) {
	return s.repo.GetByID(userID, friendID)
}

func (s *FriendService) Update(userID, friendID uuid.UUID, friend models.UpdateFriendWorkInfoInput) error {
	return s.repo.Update(userID, friendID, friend)
}

func (s *FriendService) DeleteByID(userID, friendID uuid.UUID) error {
	return s.repo.DeleteByID(userID, friendID)
}

func (s *FriendService) AddTagToFriend(friendID, tagID uuid.UUID) error {
	return s.repo.AddTagToFriend(friendID, tagID)
}

func (s *FriendService) AddTagsToFriend(userID, friendID uuid.UUID, tagIDs []models.AdditionTag) ([]uuid.UUID, error) {
	return s.repo.AddTagsToFriend(userID, friendID, tagIDs)
}

func (s *FriendService) DeleteTagFromFriend(friendID, tagID uuid.UUID) error {
	return s.repo.DeleteTagFromFriend(friendID, tagID)
}
