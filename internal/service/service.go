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

type User interface {
	Update(user models.UserUpdate, userID uuid.UUID) error
	GetByID(userID uuid.UUID) (models.User, error)
}

type Tag interface {
	Create(userID uuid.UUID, tag models.Tag) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]models.Tag, error)
	GetByID(userID, tagID uuid.UUID) (models.Tag, error)
	Update(userID, tagID uuid.UUID, tag models.Tag) error
	DeleteByID(userID, tagID uuid.UUID) error
}

type Friendlist interface {
	Create(userID uuid.UUID, friendlist models.UpdateFriendlist) (uuid.UUID, error)
	GetAll(userID uuid.UUID) ([]models.Friendlist, error)
	GetByID(userID, friendlistID uuid.UUID) (models.Friendlist, error)
	GetAllWithTags(userID uuid.UUID) ([]models.FriendlistWithTags, error)
	GetByIDWithTags(userID, friendlistID uuid.UUID) (models.FriendlistWithTags, error)
	GetAllWithFriends(userID uuid.UUID) ([]models.FriendlistWithFriends, error)
	GetByIDWithFriends(userID, friendlistID uuid.UUID) (models.FriendlistWithFriends, error)
	Update(userID, friendlistID uuid.UUID, friendlist models.UpdateFriendlist) error
	AddTagToFriendlist(friendlistID, tagID uuid.UUID) error
	DeleteTagFromFriendlist(friendlistID, tagID uuid.UUID) error
	AddFriendToFriendlist(friendlistID, friendID uuid.UUID) error
	DeleteFriendFromFriendlist(friendlistID, friendID uuid.UUID) error
	DeleteByID(userID, friendlistID uuid.UUID) error
}

type Friend interface {
	Create(userID uuid.UUID, friend models.UpdateFriendWorkInfoInput) (models.FriendIDWorkInfoID, error)
	GetAll(userID uuid.UUID) ([]models.FriendWorkInfoTags, error)
	GetByID(userID, friendID uuid.UUID) (models.FriendWorkInfoTags, error)
	Update(userID, friendID uuid.UUID, friend models.UpdateFriendWorkInfoInput) error
	DeleteByID(userID, friendID uuid.UUID) error
	AddTagToFriend(friendID, tagID uuid.UUID) error
	DeleteTagFromFriend(friendID, tagID uuid.UUID) error
}

type Event interface {
	Create(userID uuid.UUID, event models.Event) (uuid.UUID, error)
	AddFriendsToEvent(userID, eventID uuid.UUID, friendIDs []models.FriendID) ([]uuid.UUID, error)
	DeleteFriendsFromEvent(userID, eventID uuid.UUID, friendIDs []uuid.UUID) error
	GetEventsByFriendID(userID, friendID uuid.UUID) ([]models.Event, error)
	GetAll(userID uuid.UUID) ([]models.Event, error)
	GetByID(userID, eventID uuid.UUID) (models.Event, error)
	GetAllWithFriends(userID uuid.UUID) ([]models.EventWithFriends, error)
	GetByIDWithFriends(userID, eventID uuid.UUID) (models.EventWithFriends, error)
	Update(userID, eventID uuid.UUID, event models.EventUpdate) error
	DeleteByID(userID, eventID uuid.UUID) error
}

type AdditionalInfoField interface {
}

type Service struct {
	Authorization
	User
	Tag
	Friendlist
	Friend
	Event
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		User:          NewUserService(repo.User),
		Tag:           NewTagService(repo.Tag),
		Friendlist:    NewFriendlistService(repo.Friendlist),
		Friend:        NewFriendService(repo.Friend),
		Event:         NewEventService(repo.Event),
	}
}
