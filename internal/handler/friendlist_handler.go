package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

func (h *Handler) createFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.Friendlist
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	friendlistID, err := h.services.Friendlist.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"friendlist_id": friendlistID,
	})
}

func (h *Handler) getAllFriendlists(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlists, err := h.services.Friendlist.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlists": friendlists,
	})
}

func (h *Handler) getAllFriendlistsFull(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistsWithTags, err := h.services.Friendlist.GetAllWithTags(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	friendlistFull := make([]models.FriendlistFull, 0, len(friendlistsWithTags))

	friendlistsWithFriends, err := h.services.Friendlist.GetAllWithFriends(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	friends, err := h.services.Friend.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	type FriendsInFriendlist struct {
		friendlistID uuid.UUID
		friends      []models.FriendWorkInfoTags
	}

	friendsInFriendlists := make([]FriendsInFriendlist, 0, len(friendlistsWithTags))

	friendsWithTags := []models.FriendWorkInfoTags{}

	for _, friendlist := range friendlistsWithFriends {
		for _, friendFromFriendlist := range friendlist.Friends {
			for _, friend := range friends {
				if friend.Friend.ID != friendFromFriendlist.ID {
					continue
				}
				friendsWithTags = append(friendsWithTags, friend)
			}

		}
		friendsInFriendlists = append(friendsInFriendlists, FriendsInFriendlist{
			friendlistID: friendlist.Friendlist.ID,
			friends:      friendsWithTags,
		})
		friendsWithTags = nil
	}

	for _, friendlistWithTag := range friendlistsWithTags {
		for _, friendsInFriendlist := range friendsInFriendlists {
			if friendsInFriendlist.friendlistID != friendlistWithTag.Friendlist.ID {
				continue
			}
			friendlistFull = append(friendlistFull, models.FriendlistFull{
				FriendlistWithTags: friendlistWithTag,
				FriendsWithTags:    friendsInFriendlist.friends,
			})
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlists": friendlistFull,
	})
}

func (h *Handler) getFriendlistByID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friendlist, err := h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlist": friendlist,
	})
}

func (h *Handler) getFriendlistByIDFull(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friendlistWithTags, err := h.services.Friendlist.GetByIDWithTags(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	friendlistWithFriends, err := h.services.Friendlist.GetByIDWithFriends(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	friends := []models.FriendWorkInfoTags{}
	for _, friend := range friendlistWithFriends.Friends {
		friendWithTags, err := h.services.Friend.GetByID(userID, friend.ID)
		if err != nil {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		friends = append(friends, friendWithTags)
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlist_full": models.FriendlistFull{
			FriendlistWithTags: friendlistWithTags,
			FriendsWithTags:    friends,
		},
	})
}

func (h *Handler) getAllFriendlistsWithTags(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlists, err := h.services.Friendlist.GetAllWithTags(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlists": friendlists,
	})
}

func (h *Handler) getFriendlistByIDWithTags(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friendlist, err := h.services.Friendlist.GetByIDWithTags(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlist": friendlist,
	})
}

func (h *Handler) getAllFriendlistsWithFriends(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlists, err := h.services.Friendlist.GetAllWithFriends(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlists": friendlists,
	})
}

func (h *Handler) getFriendlistByIDWithFriends(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friendlist, err := h.services.Friendlist.GetByIDWithFriends(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friendlist": friendlist,
	})
}

func (h *Handler) updateFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	var payload models.Friendlist
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.Update(userID, friendlistID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found or already deleted: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.DeleteByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) addTagToFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}

	var payload models.AdditionTag
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Tag.GetByID(userID, payload.TagID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("tag not found: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.AddTagToFriendlist(friendlistID, payload.TagID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, statusResponse{
		Status: "ok",
	},
	)
}

func (h *Handler) addFriendToFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}

	var payload models.AdditionFriendToFriendlist
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Friend.GetByID(userID, payload.FriendID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friend not found: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.AddFriendToFriendlist(friendlistID, payload.FriendID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, statusResponse{
		Status: "ok",
	},
	)
}

func (h *Handler) deleteTagFromFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.DeleteTagFromFriendlist(friendlistID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteFriendFromFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friendID, err := uuid.Parse(c.Param("friend_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}

	err = h.services.Friendlist.DeleteFriendFromFriendlist(friendlistID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
