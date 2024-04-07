package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{
		repo: repo,
	}
}

func (s *TagService) Create(userID uuid.UUID, tag models.Tag) (uuid.UUID, error) {
	return s.repo.Create(userID, tag)
}

func (s *TagService) GetAll(userID uuid.UUID) ([]models.Tag, error) {
	return s.repo.GetAll(userID)
}

func (s *TagService) GetByID(userID, tagID uuid.UUID) (models.Tag, error) {
	return s.repo.GetByID(userID, tagID)
}

func (s *TagService) Update(userID, tagID uuid.UUID, tag models.Tag) error {
	return s.repo.Update(userID, tagID, tag)
}

func (s *TagService) DeleteByID(userID, tagID uuid.UUID) error {
	return s.repo.DeleteByID(userID, tagID)
}
