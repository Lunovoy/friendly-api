package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Create Tag
// @Security ApiKeyAuth
// @Tags tag
// @Description create tag
// @ID create-tag
// @Accept  json
// @Produce  json
// @Param input body models.Tag true "Tag info"
// @Success 201 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tag [post]
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
		if err.Error() == "tag already exists" {
			c.JSON(http.StatusOK, map[string]any{
				"message": err.Error(),
				"tag_id":  tagID,
			})
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, map[string]any{
		"tag_id": tagID,
	})

}

// @Summary Get All Tags
// @Security ApiKeyAuth
// @Tags tag
// @Description get all tags
// @ID get-all-tags
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllTagsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tag [get]
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

	c.JSON(http.StatusOK, getAllTagsResponse{
		Data: tags,
	})

}

// @Summary Get Tag By Id
// @Security ApiKeyAuth
// @Tags tag
// @Description get tag by id
// @ID get-tag-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Tag id"
// @Success 200 {object} models.Tag
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tag/{id} [get]
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

// @Summary Update Tag
// @Security ApiKeyAuth
// @Tags tag
// @Description update tag
// @ID update-tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag id"
// @Param input body models.Tag true "Tag info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tag/{id} [put]
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

// @Summary Delete Tag
// @Security ApiKeyAuth
// @Tags tag
// @Description delete tag
// @ID delete-tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tag/{id} [delete]
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
