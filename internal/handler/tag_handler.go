package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

func (h *Handler) createTag(c *gin.Context) {

	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.Tag
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tagID, err := h.services.Tag.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"tag_id": tagID,
	})

}
func (h *Handler) getAllTags(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	tags, err := h.services.Tag.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"tags": tags,
	})

}
func (h *Handler) getTagByID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	tag, err := h.services.Tag.GetByID(userID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"tag": tag,
	})
}
func (h *Handler) updateTag(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	var payload models.Tag
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Tag.GetByID(userID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("tag not found: %s", err.Error()))
		return
	}

	err = h.services.Tag.Update(userID, tagID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
func (h *Handler) deleteTag(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Tag.GetByID(userID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("tag not found or already deleted: %s", err.Error()))
		return
	}

	err = h.services.Tag.DeleteByID(userID, tagID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
