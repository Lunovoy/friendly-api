package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/models"
)

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
