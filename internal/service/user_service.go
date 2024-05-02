package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByID(userID uuid.UUID) (models.User, error) {
	return s.repo.GetByID(userID)
}

func (s *UserService) Update(user models.UserUpdate, userID uuid.UUID) error {
	return s.repo.Update(user, userID)
}
