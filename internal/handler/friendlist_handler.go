package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Create Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description create friendlist
// @ID create-friendlist
// @Accept  json
// @Produce  json
// @Param input body models.UpdateFriendlist true "Friendlist info"
// @Success 201 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist [post]
func (h *Handler) createFriendlist(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.UpdateFriendlist
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

// @Summary Get All Friendlists
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get all friendlists
// @ID get-all-friendlists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllFriendlistsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist [get]
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

	c.JSON(http.StatusOK, getAllFriendlistsResponse{
		Data: friendlists,
	})
}

// @Summary Get All Friendlists Full
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get all friendlists full
// @ID get-all-friendlists-full
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllFriendlistsFullResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/full [get]
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

	type friendsInFriendlist struct {
		friendlistID uuid.UUID
		friends      []models.FriendWorkInfoTags
	}

	friendsInFriendlists := make([]friendsInFriendlist, 0, len(friendlistsWithTags))

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
		friendsInFriendlists = append(friendsInFriendlists, friendsInFriendlist{
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

	c.JSON(http.StatusOK, getAllFriendlistsFullResponse{
		Data: friendlistFull,
	})
}

// @Summary Get Friendlist By Id
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get friendlist by id
// @ID get-friendlist-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Success 200 {object} models.Friendlist
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id} [get]
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

// @Summary Get Friendlist Full By ID
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get friendlist full by id
// @ID get-friendlist-full-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Success 200 {object} models.FriendlistFull
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/full [get]
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

// @Summary Get All Friendlists With Tags
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get all friendlists with tags
// @ID get-all-friendlists-with-tags
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllFriendlistsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/tag [get]
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

	c.JSON(http.StatusOK, getAllFriendlistsWithTagsResponse{
		Data: friendlists,
	})
}

// @Summary Get Friendlist By ID With Tags
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get friendlist by id with tags
// @ID get-friendlist-by-id-with-tags
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Success 200 {object} models.FriendlistWithTags
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/tag [get]
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

// @Summary Get All Friendlists With Friends
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get all friendlists with friends
// @ID get-all-friendlists-with-friends
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllFriendlistsWithFriendsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/friend [get]
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

	c.JSON(http.StatusOK, getAllFriendlistsWithFriendsResponse{
		Data: friendlists,
	})
}

// @Summary Get Friendlist By ID With Friends
// @Security ApiKeyAuth
// @Tags friendlist
// @Description get friendlist by id with friends
// @ID get-friendlist-by-id-with-friends
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Success 200 {object} models.FriendlistWithFriends
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/friend [get]
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

// @Summary Update Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description update friendlist
// @ID update-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Param input body models.UpdateFriendlist true "Friendlist info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id} [put]
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

	var payload models.UpdateFriendlist
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	friendlist, err := h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found: %s", err.Error()))
		return
	}
	oldImageID := friendlist.ImageID

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

	if payload.ImageID != nil && oldImageID != uuid.Nil && *payload.ImageID != oldImageID {
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

// @Summary Delete Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description delete friendlist
// @ID delete-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id} [delete]
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

	friendlist, err := h.services.Friendlist.GetByID(userID, friendlistID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("friendlist not found or already deleted: %s", err.Error()))
		return
	}

	oldImageID := friendlist.ImageID

	err = h.services.Friendlist.DeleteByID(userID, friendlistID)
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

// @Summary Add Tag To Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description add tag to friendlist
// @ID add-tag-to-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Param input body models.AdditionTag true "Add Tag"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/tag [post]
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

// @Summary Add Friend To Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description add friend to friendlist
// @ID add-friend-to-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Param input body models.AdditionFriendToFriendlist true "Add Friend"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/friend [post]
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

// @Summary Delete Tag From Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description delete tag from friendlist
// @ID delete-tag-from-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Param tag_id path string true "Tag id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/tag/{tag_id} [delete]
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

// @Summary Delete Friend From Friendlist
// @Security ApiKeyAuth
// @Tags friendlist
// @Description delete friend from friendlist
// @ID delete-friend-from-friendlist
// @Accept  json
// @Produce  json
// @Param id path string true "Friendlist id"
// @Param friend_id path string true "Friend id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/friendlist/{id}/friend/{friend_id} [delete]
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
