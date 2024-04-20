package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lunovoy/friendly/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (uuid.UUID, error)
	GetUserByMail(mail string) (models.User, error)
	GetUserByID(id uuid.UUID) (models.User, error)
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
	GetAllWithFriends(userID uuid.UUID) ([]models.FriendlistWithFriends, error)
	GetByIDWithFriends(userID, friendlistID uuid.UUID) (models.FriendlistWithFriends, error)
	Update(userID, friendlistID uuid.UUID, friendlist models.Friendlist) error
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

type AdditionalInfoField interface {
}

type Repository struct {
	Authorization
	Tag
	Friendlist
	Friend
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Tag:           NewTagPostgres(db),
		Friendlist:    NewFriendlistPostgres(db),
		Friend:        NewFriendPostgres(db),
	}
}
