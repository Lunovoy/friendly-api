package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Create Event
// @Security ApiKeyAuth
// @Tags event
// @Description create event
// @ID create-event
// @Accept  json
// @Produce  json
// @Param input body models.Event true "Event info"
// @Success 201 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event [post]
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

// TODO: add friends addition while creating event
// add deletion friends from events

func (h *Handler) addFriendsToEvent(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload []models.FriendID
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	ids, err := h.services.Event.AddFriendsToEvent(userID, eventID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"status": "ok",
		"IDs":    ids,
	})
}

// @Summary Get Event By Friend ID
// @Security ApiKeyAuth
// @Tags event
// @Description get event by friend id
// @ID get-event-by-friend-id
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllEventsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event//friend/:friend_id [get]
func (h *Handler) getEventsByFriendID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	friendID, err := uuid.Parse(c.Param("friend_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	events, err := h.services.Event.GetEventsByFriendID(userID, friendID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllEventsResponse{
		Data: events,
	})
}

// @Summary Get All Events
// @Security ApiKeyAuth
// @Tags event
// @Description get all events
// @ID get-all-events
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllEventsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event [get]
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

	c.JSON(http.StatusOK, getAllEventsResponse{
		Data: events,
	})

}

// @Summary Get Event By Id
// @Security ApiKeyAuth
// @Tags event
// @Description get event by id
// @ID get-event-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Event
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [get]
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

// @Summary Get All Events With Friends
// @Security ApiKeyAuth
// @Tags event
// @Description get all events with friends
// @ID get-all-events-with-friends
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllEventsWithFriendsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/friends [get]
func (h *Handler) getAllEventsWithFriends(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	events, err := h.services.Event.GetAllWithFriends(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllEventsWithFriendsResponse{
		Data: events,
	})

}

// @Summary Get Event By Id With Friends
// @Security ApiKeyAuth
// @Tags event
// @Description get event by id with friends
// @ID get-event-by-id-with-friends
// @Accept  json
// @Produce  json
// @Success 200 {object} models.EventWithFriends
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id/friends [get]
func (h *Handler) getEventByIDWithFriends(c *gin.Context) {
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

	event, err := h.services.Event.GetByIDWithFriends(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"event": event,
	})
}

// @Summary Update Event
// @Security ApiKeyAuth
// @Tags event
// @Description update event
// @ID update-event
// @Accept  json
// @Produce  json
// @Param input body models.EventUpdate true "Event info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [put]
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

	var payload models.EventUpdate
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

// @Summary Delete Event
// @Security ApiKeyAuth
// @Tags event
// @Description delete event
// @ID delete-event
// @Accept  json
// @Produce  json
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id [delete]
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
