package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/service"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	maxFileSize    = 5 << 20 // 5MB
	uploadDir      = "./images/"
	imageExtension = ".jpg"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api")
	{

		auth := api.Group("/auth")
		{
			auth.POST("/sign-in", h.signIn)
			auth.POST("/sign-up", h.signUp)
		}

		profile := api.Group("/profile", h.userIdentity)
		{
			profile.GET("/", h.getProfile)
			profile.PUT("/", h.updateProfile)
		}

		tags := api.Group("/tag", h.userIdentity)
		{
			tags.POST("/", h.createTag)
			tags.GET("/", h.getAllTags)
			tags.GET("/:id", h.getTagByID)
			tags.PUT("/:id", h.updateTag)
			tags.DELETE("/:id", h.deleteTag)
		}

		friendlist := api.Group("/friendlist", h.userIdentity)
		{
			friendlist.POST("/", h.createFriendlist)
			friendlist.GET("/", h.getAllFriendlists)
			friendlist.GET("/full", h.getAllFriendlistsFull)
			friendlist.GET("/:id", h.getFriendlistByID)
			friendlist.GET("/:id/full", h.getFriendlistByIDFull)
			friendlist.GET("/tag", h.getAllFriendlistsWithTags)
			friendlist.GET("/:id/tag", h.getFriendlistByIDWithTags)
			friendlist.GET("/friend", h.getAllFriendlistsWithFriends)
			friendlist.GET("/:id/friend", h.getFriendlistByIDWithFriends)
			friendlist.PUT("/:id", h.updateFriendlist)
			friendlist.DELETE("/:id", h.deleteFriendlist)
			friendlist.POST("/:id/tag", h.addTagToFriendlist)
			friendlist.POST("/:id/friend", h.addFriendToFriendlist)
			friendlist.DELETE("/:id/tag/:tag_id", h.deleteTagFromFriendlist)
			friendlist.DELETE("/:id/friend/:friend_id", h.deleteFriendFromFriendlist)
		}

		friend := api.Group("/friend", h.userIdentity)
		{
			friend.POST("/", h.createFriend)
			friend.GET("/", h.getAllFriends)
			friend.GET("/:id", h.getFriendByID)
			friend.PUT("/:id", h.updateFriend)
			friend.DELETE("/:id", h.deleteFriend)
			friend.POST("/:id/tag", h.addTagToFriend)
			friend.DELETE("/:id/tag/:tag_id", h.deleteTagFromFriend)
		}

		event := api.Group("/event", h.userIdentity)
		{
			event.POST("/", h.createEvent)
			event.POST("/:id/friends", h.addFriendsToEvent)
			event.GET("/friend/:friend_id", h.getEventsByFriendID)
			event.GET("/", h.getAllEvents)
			event.GET("/full", h.getAllEventsFull)
			event.GET("/:id", h.getEventByID)
			event.GET("/friends", h.getAllEventsWithFriends)
			event.GET("/:id/full", h.getEventByIDFull)
			event.PUT("/:id", h.updateEvent)
			event.DELETE("/:id", h.deleteEvent)
		}

		reminder := api.Group("/reminder", h.userIdentity)
		{
			reminder.POST("/", h.createReminder)
			reminder.GET("/", h.getAllReminders)
			reminder.GET("/event/:event_id", h.getAllRemindersByEventID)
			reminder.GET("/:id", h.getReminderByID)
			reminder.PUT("/:id", h.updateReminder)
			reminder.DELETE("/:id", h.deleteReminder)
		}

		additionalInfoField := api.Group("/additional-field", h.userIdentity)
		{
			additionalInfoField.POST("/", h.createAdditionalInfoField)
			additionalInfoField.GET("/", h.getAllAdditionalFields)
			additionalInfoField.GET("/:id", h.getAdditionalFieldByID)
			additionalInfoField.PUT("/:id", h.updateAdditionalField)
			additionalInfoField.DELETE("/:id", h.deleteAdditionalField)
		}

		image := api.Group("/image", h.userIdentity)
		{
			image.POST("/", h.uploadImage)
			image.GET("/:id/:res", h.getImage)
			image.DELETE("/:id", h.deleteImage)
		}
	}

	return router
}
