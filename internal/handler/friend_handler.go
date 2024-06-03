package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Create Friend
// @Security ApiKeyAuth
// @Tags friend
// @Description create friend
// @ID create-friend
// @Accept  json
// @Produce  json
// @Param input body models.UpdateFriendWorkInfoTagsInput true "Friend info"
// @Success 201 {object} any
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend [post]
func (h *Handler) createFriend(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.UpdateFriendWorkInfoTagsInput
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	friendPayload := models.UpdateFriendWorkInfoInput{
		Friend:   payload.Friend,
		WorkInfo: payload.WorkInfo,
	}

	friendIDWorkID, err := h.services.Friend.Create(userID, friendPayload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.TagIDs != nil && len(payload.TagIDs) != 0 {
		_, err = h.services.Friend.AddTagsToFriend(userID, friendIDWorkID.FriendID, payload.TagIDs)
		if err != nil {
			if delErr := h.services.Friend.DeleteByID(userID, friendIDWorkID.FriendID); delErr != nil {
				newErrorResponse(c, http.StatusInternalServerError, delErr.Error())
			}
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	friend, err := h.services.Friend.GetByID(userID, friendIDWorkID.FriendID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if friend.Friend.DOB.Valid {

		eventDate := time.Date(time.Now().Year(), friend.Friend.DOB.Time.Month(), friend.Friend.DOB.Time.Day(), friend.Friend.DOB.Time.Hour(), friend.Friend.DOB.Time.Minute(), 0, 0, friend.Friend.DOB.Time.Location())
		if eventDate.Before(time.Now()) {
			eventDate = eventDate.AddDate(1, 0, 0)
		}

		eventID, err := h.services.Event.Create(userID, models.Event{
			Title:       fmt.Sprintf("День рождение: %s %s", friend.Friend.FirstName, friend.Friend.LastName),
			Description: fmt.Sprintf("%s %s", friend.WorkInfo.City, friend.WorkInfo.Company),
			Frequency:   "annually",
			StartDate: sql.NullTime{
				Time:  eventDate,
				Valid: true,
			},
			EndDate: sql.NullTime{
				Time:  eventDate.Add(5 * time.Minute),
				Valid: true,
			},
		})
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error while creating dob event of friend: %s", err.Error()))
			return
		}
		log.Printf("event dob created: %s", eventID)
		var friends []models.FriendID
		friends = append(friends, models.FriendID{FriendID: friend.Friend.ID})
		_, err = h.services.Event.AddFriendsToEvent(userID, eventID, friends)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error while adding friend to event of friend dob: %s", err.Error()))
			return
		}
		_, err = h.services.Reminder.Create(userID, models.Reminder{EventID: eventID, MinutesUntilEvent: 0})
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error while adding reminder to event of friend dob: %s", err.Error()))
			return
		}
		friends = nil
	}

	c.JSON(http.StatusCreated, map[string]any{
		"friend_id":    friendIDWorkID.FriendID,
		"work_info_id": friendIDWorkID.WorkInfoID,
	})
}

// @Summary Get All Friends
// @Security ApiKeyAuth
// @Tags friend
// @Description get all friends
// @ID get-all-friends
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllFriendsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend [get]
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

	c.JSON(http.StatusOK, getAllFriendsResponse{
		Data: friends,
	})
}

// @Summary Get Friend By ID
// @Security ApiKeyAuth
// @Tags friend
// @Description get friend by id
// @ID get-friend-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Friend id"
// @Success 200 {object} models.FriendWorkInfoTags
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend/{id} [get]
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

// @Summary Update Friend
// @Security ApiKeyAuth
// @Tags friend
// @Description update friend
// @ID update-friend
// @Accept  json
// @Produce  json
// @Param id path string true "Friend id"
// @Param input body models.UpdateFriendWorkInfoInput true "Friend info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend/{id} [put]
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

// @Summary Delete Friend
// @Security ApiKeyAuth
// @Tags friend
// @Description delete friend
// @ID delete-friend
// @Accept  json
// @Produce  json
// @Param id path string true "Friend id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend/{id} [delete]
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

// @Summary Add Tag To Friend
// @Security ApiKeyAuth
// @Tags friend
// @Description add tag to friend
// @ID add-tag-to-friend
// @Accept  json
// @Produce  json
// @Param id path string true "Friend id"
// @Param input body models.AdditionTag true "Add Tag"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend/{id}/tag [post]
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

// @Summary Delete Tag From Friend
// @Security ApiKeyAuth
// @Tags friend
// @Description delete tag from friend
// @ID delete-tag-from-friend
// @Accept  json
// @Produce  json
// @Param id path string true "Friend id"
// @Param tag_id path string true "Tag id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friend/{id}/tag/{tag_id} [delete]
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
