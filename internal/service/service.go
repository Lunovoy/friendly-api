package service

import (
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/lunovoy/friendly/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (uuid.UUID, error)
	GetUserByMail(mail, password string) (models.User, error)
	GenerateToken(userID uuid.UUID) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
}

type Tag interface {
	Create(userID uuid.UUID, tag models.Tag) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]models.Tag, error)
	GetByID(userID, tagID uuid.UUID) (models.Tag, error)
	Update(userID, tagID uuid.UUID, tag models.Tag) error
	DeleteByID(userID, tagID uuid.UUID) error
}

type Friendlist interface {
	Create(userID uuid.UUID, friendlist models.Friendlist) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]models.Friendlist, error)
	GetByID(userID, friendlistID uuid.UUID) (models.Friendlist, error)
	GetAllWithTags(userID uuid.UUID) ([]models.FriendlistWithTags, error)
	GetByIDWithTags(userID, friendlistID uuid.UUID) (models.FriendlistWithTags, error)
	Update(userID, friendlistID uuid.UUID, friendlist models.Friendlist) error
	AddTagToFriendlist(friendlistID, tagID uuid.UUID) error
	DeleteTagFromFriendlist(friendlistID, tagID uuid.UUID) error
	AddFriendToFriendlist(friendlistID, friendID uuid.UUID) error
	DeleteFriendFromFriendlist(friendlistID, friendID uuid.UUID) error
	DeleteByID(userID, friendlistID uuid.UUID) error
}

type Friend interface {
	Create(userID uuid.UUID, friend models.Friend, workInfo models.WorkInfo) (models.FriendIDWorkInfoID, error)
	GetAll(userID uuid.UUID) ([]models.FriendWorkInfo, error)
	GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfo, error)
	Update(userID, friendID uuid.UUID, friend models.Friend) error
	DeleteByID(userID, friendID uuid.UUID) error
}

type AdditionalInfoField interface {
}

type Service struct {
	Authorization
	Tag
	Friendlist
	Friend
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Tag:           NewTagService(repo.Tag),
		Friendlist:    NewFriendlistService(repo.Friendlist),
		Friend:        NewFriendService(repo.Friend),
	}
}
