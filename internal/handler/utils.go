package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"

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
func saveFile(c *gin.Context, file *multipart.FileHeader, filePath string) error {

	err := c.SaveUploadedFile(file, filePath)
	return err
}
func deleteFile(filePath string) error {

	err := os.Remove(filePath)

	return err

}
