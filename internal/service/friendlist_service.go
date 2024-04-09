package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type FriendlistService struct {
	repo repository.Friendlist
}

func NewFriendlistService(repo repository.Friendlist) *FriendlistService {
	return &FriendlistService{
		repo: repo,
	}
}

func (s *FriendlistService) Create(userID uuid.UUID, friendlist models.Friendlist) (uuid.UUID, error) {
	return s.repo.Create(userID, friendlist)
}

func (s *FriendlistService) GetAll(userID uuid.UUID) ([]models.Friendlist, error) {
	return s.repo.GetAll(userID)
}

func (s *FriendlistService) GetByID(userID, friendlistID uuid.UUID) (models.Friendlist, error) {
	return s.repo.GetByID(userID, friendlistID)
}

func (s *FriendlistService) GetAllWithTags(userID uuid.UUID) ([]models.FriendlistWithTags, error) {
	return s.repo.GetAllWithTags(userID)
}

func (s *FriendlistService) GetByIDWithTags(userID, friendlistID uuid.UUID) (models.FriendlistWithTags, error) {
	return s.repo.GetByIDWithTags(userID, friendlistID)
}

func (s *FriendlistService) GetAllWithFriends(userID uuid.UUID) ([]models.FriendlistWithFriends, error) {
	return s.repo.GetAllWithFriends(userID)
}

func (s *FriendlistService) GetByIDWithFriends(userID, friendlistID uuid.UUID) (models.FriendlistWithFriends, error) {
	return s.repo.GetByIDWithFriends(userID, friendlistID)
}

func (s *FriendlistService) Update(userID, friendlistID uuid.UUID, friendlist models.Friendlist) error {
	return s.repo.Update(userID, friendlistID, friendlist)
}

func (s *FriendlistService) AddTagToFriendlist(friendlistID, tagID uuid.UUID) error {
	return s.repo.AddTagToFriendlist(friendlistID, tagID)
}

func (s *FriendlistService) AddFriendToFriendlist(friendlistID, friendID uuid.UUID) error {
	return s.repo.AddFriendToFriendlist(friendlistID, friendID)
}

func (s *FriendlistService) DeleteTagFromFriendlist(friendlistID, tagID uuid.UUID) error {
	return s.repo.DeleteTagFromFriendlist(friendlistID, tagID)
}
func (s *FriendlistService) DeleteFriendFromFriendlist(friendlistID, friendID uuid.UUID) error {
	return s.repo.DeleteFriendFromFriendlist(friendlistID, friendID)
}

func (s *FriendlistService) DeleteByID(userID, friendlistID uuid.UUID) error {
	return s.repo.DeleteByID(userID, friendlistID)
}
