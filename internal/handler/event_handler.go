package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

func (h *Handler) createEvent(c *gin.Context) {

	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.Event
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	eventID, err := h.services.Event.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"event_id": eventID,
	})

}

func (h *Handler) addFriendsToEvent(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.Event
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	eventID, err := h.services.Event.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
	})
}

func (h *Handler) getEventsByFriendID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}
}

func (h *Handler) getAllEvents(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	events, err := h.services.Event.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"events": events,
	})

}

func (h *Handler) getEventByID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	event, err := h.services.Event.GetByID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"event": event,
	})
}

func (h *Handler) updateEvent(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	var payload models.Event
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Event.GetByID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("event not found: %s", err.Error()))
		return
	}

	err = h.services.Event.Update(userID, eventID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteEvent(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Event.GetByID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("event not found or already deleted: %s", err.Error()))
		return
	}

	err = h.services.Event.DeleteByID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
