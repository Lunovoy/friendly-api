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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG files are allowed."})
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
