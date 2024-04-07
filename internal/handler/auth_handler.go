package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunovoy/friendly/internal/models"
)

type signInPayload struct {
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var payload models.User

	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.services.Authorization.CreateUser(payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"user_id": userID,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var payload signInPayload

	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.Authorization.GetUserByMail(payload.Mail, payload.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "not found, invalid mail or password")
		return
	}

	token, err := h.services.Authorization.GenerateToken(user.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}
