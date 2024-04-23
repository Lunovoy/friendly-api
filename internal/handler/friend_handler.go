package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

func (h *Handler) createFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.UpdateFriendWorkInfoInput
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	friendIDWorkID, err := h.services.Friend.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"friend_id":    friendIDWorkID.FriendID,
		"work_info_id": friendIDWorkID.WorkInfoID,
	})
}

func (h *Handler) getAllFriends(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friends, err := h.services.Friend.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friends": friends,
	})
}
func (h *Handler) getFriendByID(c *gin.Context) {
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

	friend, err := h.services.Friend.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"friend": friend,
	})
}

func (h *Handler) updateFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	var payload models.UpdateFriendWorkInfoInput
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	friend, err := h.services.Friend.GetByID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friend not found: %s", err.Error()))
		return
	}
	oldImageID := friend.Friend.ImageID

	err = h.services.Friend.Update(userID, friendID, payload)
	if err != nil {
		if payload.Friend.ImageID != nil {
			deleteFile(fmt.Sprintf("%s%s%s", uploadDir, payload.Friend.ImageID, imageExtension))
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.Friend.ImageID != nil && oldImageID != uuid.Nil && *payload.Friend.ImageID != oldImageID {
		err := deleteFile(fmt.Sprintf("%s%s%s", uploadDir, oldImageID, imageExtension))
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	friend, err := h.services.Friend.GetByID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friend not found or already deleted: %s", err.Error()))
		return
	}

	oldImageID := friend.Friend.ImageID

	err = h.services.Friend.DeleteByID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if oldImageID != uuid.Nil {
		err := deleteFile(fmt.Sprintf("%s%s%s", uploadDir, oldImageID, imageExtension))
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) addTagToFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friend.GetByID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friend not found: %s", err.Error()))
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

	err = h.services.Friend.AddTagToFriend(friendID, payload.TagID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, statusResponse{
		Status: "ok",
	},
	)
}

func (h *Handler) deleteTagFromFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Friend.GetByID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friend not found: %s", err.Error()))
		return
	}

	err = h.services.Friend.DeleteTagFromFriend(friendID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
