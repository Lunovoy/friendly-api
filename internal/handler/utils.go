package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getUserIDFromCtx(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return uuid.Nil, errors.New("user id not found")
	}

	convertedID, ok := id.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id invalid type")
		return uuid.Nil, errors.New("user id invalid type")
	}

	return convertedID, nil
}
