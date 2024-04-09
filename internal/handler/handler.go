package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

	api := router.Group("/api", h.userIdentity)
	{
		tags := api.Group("/tag")
		{
			tags.POST("/", h.createTag)
			tags.GET("/", h.getAllTags)
			tags.GET("/:id", h.getTagByID)
			tags.PUT("/:id", h.updateTag)
			tags.DELETE("/:id", h.deleteTag)
		}

		friendlist := api.Group("/friendlist")
		{
			friendlist.POST("/", h.createFriendlist)
			friendlist.GET("/", h.getAllFriendlists)
			friendlist.GET("/:id", h.getFriendlistByID)
			friendlist.GET("/tag", h.getAllFriendlistsWithTags)
			friendlist.GET("/:id/tag", h.getFriendlistByIDWithTags)
			friendlist.PUT("/:id", h.updateFriendlist)
			friendlist.DELETE("/:id", h.deleteFriendlist)
			friendlist.POST("/:id/tag", h.addTagToFriendlist)
			friendlist.POST("/:id/friend", h.addFriendToFriendlist)
			friendlist.DELETE("/:id/tag/:tag_id", h.deleteTagFromFriendlist)
			friendlist.DELETE("/:id/friend/:friend_id", h.deleteFriendFromFriendlist)
		}

		friend := api.Group("/friend")
		{
			friend.POST("/", h.createFriend)
			friend.GET("/", h.getAllFriends)
			friend.GET("/:id", h.getFriendByID)
			friend.PUT("/:id", h.updateFriend)
			friend.DELETE("/:id", h.deleteFriend)
		}
	}
	return router
}
