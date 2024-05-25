package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

// @Summary Upload Image
// @Security ApiKeyAuth
// @Tags image
// @Description upload image
// @ID upload-image
// @Accept  multipart/form-data
// @Produce  json
// @Param image formData file true "Image file to upload"
// @Success 201 {object} map[string]any "Successfully uploaded image"
// @Failure 400 {object} errorResponse
// @Failure 415 {object} errorResponse "Invalid file type. Only JPG files are allowed."
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure default {object} errorResponse
// @Router /api/image [post]
func (h *Handler) uploadImage(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the maximum allowed size 5MB"})
	}
	uuid := uuid.NewString()

	file, err := c.FormFile("image")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// Проверяем тип файла
	if file.Header.Get("Content-Type") != "image/jpeg" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Invalid file type. Only JPG files are allowed."})
		return
	}

	fileExt := strings.Split(file.Filename, ".")[1]
	filename := fmt.Sprintf("%s.%s", uuid, fileExt)
	filePath := uploadDir + filename

	// Проверяем, существует ли файл с таким именем, и если существует - удаляем его
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete existing file",
			})
			return
		}
	}

	// Сохраняем файл на сервере

	if err := saveFile(c, file, filePath); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error saving file: %s", err.Error()))
		return
	}

	data := map[string]interface{}{
		"header":   file.Header,
		"image_id": uuid,
		"size":     file.Size,
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

// @Summary Get Image
// @Security ApiKeyAuth
// @Tags image
// @Description get image by ID and resolution
// @ID get-image
// @Accept  json
// @Produce  image/jpeg
// @Param id path string true "Image ID"
// @Param res path int true "Resolution" Format(int64)
// @Success 200 {string} image/jpeg "Successfully retrieved image"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse "Image not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure default {object} errorResponse
// @Router /api/image/{id}/{res} [get]
func (h *Handler) getImage(c *gin.Context) {
	imageID := c.Param("id")
	resolutionString := c.Param("res")
	resolution, err := strconv.Atoi(resolutionString)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid param")
		return
	}
	filePath := fmt.Sprintf("%s%s%s", uploadDir, imageID, imageExtension)

	// Проверяем, существует ли файл с таким именем
	_, err = os.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Декодируем изображение
	img, _, err := image.Decode(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	// Изменяем размер изображения
	img = resize.Resize(uint(resolution), uint(resolution), img, resize.Lanczos3)

	// Устанавливаем заголовок Content-Type для изображения
	c.Header("Content-Type", "image/jpeg")

	// Копируем содержимое файла в ответ
	if err := jpeg.Encode(c.Writer, img, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file contents"})
		return
	}

}

// @Summary Delete Image
// @Security ApiKeyAuth
// @Tags image
// @Description delete image by ID
// @ID delete-image
// @Accept  json
// @Produce  json
// @Param id path string true "Image ID"
// @Success 200 {object} map[string]string "Image deleted successfully"
// @Failure 404 {object} errorResponse "Image not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure default {object} errorResponse
// @Router /api/image/{id} [delete]
func (h *Handler) deleteImage(c *gin.Context) {
	friendID := c.Param("id")
	filePath := fmt.Sprintf("%s%s", uploadDir, friendID+".jpg")

	// Удаляем файл с сервера
	if err := deleteFile(filePath); err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Failed to delete file: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
