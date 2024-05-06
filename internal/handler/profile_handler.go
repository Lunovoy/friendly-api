package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Get Profile
// @Security ApiKeyAuth
// @Tags profile
// @Description get profile
// @ID get-profile
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/profile/ [get]
func (h *Handler) getProfile(c *gin.Context) {

	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	user, err := h.services.User.GetByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"profile": user,
	})

}

// @Summary Update Profile
// @Security ApiKeyAuth
// @Tags profile
// @Description update user profile
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body models.UserUpdate true "User update payload"
// @Success 200 {object} statusResponse "Successfully updated profile"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 404 {object} errorResponse "User not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure default {object} errorResponse
// @Router /api/profile [put]
func (h *Handler) updateProfile(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.UserUpdate
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.User.GetByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("user not found: %s", err.Error()))
		return
	}

	err = h.services.User.Update(payload, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
