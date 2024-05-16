package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/models"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getAllEventsResponse struct {
	Data []models.Event `json:"data"`
}

type getAllEventsFullInfo struct {
	Data []*models.EventWithFriendsAndReminders `json:"data"`
}

type getAllEventsWithFriendsResponse struct {
	Data []models.EventWithFriends `json:"data"`
}

type getAllTagsResponse struct {
	Data []models.Tag `json:"data"`
}

type getAllFriendsResponse struct {
	Data []models.FriendWorkInfoTags `json:"data"`
}

type getAllFriendlistsResponse struct {
	Data []models.Friendlist `json:"data"`
}

type getAllFriendlistsFullResponse struct {
	Data []models.FriendlistFull `json:"data"`
}

type getAllFriendlistsWithTagsResponse struct {
	Data []models.FriendlistWithTags `json:"data"`
}

type getAllFriendlistsWithFriendsResponse struct {
	Data []models.FriendlistWithFriends `json:"data"`
}

type getAllRemindersResponse struct {
	Data []models.Reminder `json:"data"`
}

type statusResponse struct {
	Status string
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
