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
// @Param input body models.EventWithFriendIDsAndReminders true "Event info with friends and reminders"
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

	var payload models.EventWithFriendIDsAndReminders
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	eventID, err := h.services.Event.Create(userID, payload.Event)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.Reminders != nil && len(payload.Reminders) != 0 {
		_, err = h.services.Reminder.CreateBulk(userID, eventID, payload.Reminders)
		if err != nil {
			if delErr := h.services.Event.DeleteByID(userID, eventID); delErr != nil {
				newErrorResponse(c, http.StatusInternalServerError, delErr.Error())
			}
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if payload.FriendIDs != nil && len(payload.FriendIDs) != 0 {
		_, err = h.services.Event.AddFriendsToEvent(userID, eventID, payload.FriendIDs)
		if err != nil {
			if delErr := h.services.Event.DeleteByID(userID, eventID); delErr != nil {
				newErrorResponse(c, http.StatusInternalServerError, delErr.Error())
			}
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
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

// @Summary Get Event By Id With Friends And Reminders
// @Security ApiKeyAuth
// @Tags event
// @Description get event by id with friends and reminders
// @ID get-event-by-id-with-friends-and-reminders
// @Accept  json
// @Produce  json
// @Success 200 {object} models.EventWithFriendsAndReminders
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/:id/full [get]
func (h *Handler) getEventByIDFull(c *gin.Context) {
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

	eventWithFriends, err := h.services.Event.GetByIDWithFriends(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	reminders, err := h.services.Reminder.GetAllByEventID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	eventFull := models.EventWithFriendsAndReminders{
		Event:     eventWithFriends.Event,
		Friends:   eventWithFriends.Friends,
		Reminders: reminders,
	}

	c.JSON(http.StatusOK, map[string]any{
		"event": eventFull,
	})
}

// @Summary Get All Events Full Info
// @Security ApiKeyAuth
// @Tags event
// @Description get all events full info
// @ID get-all-events-full-info
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllEventsFullInfo
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/event/full [get]
func (h *Handler) getAllEventsFull(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	eventsWithFriends, err := h.services.Event.GetAllWithFriends(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	reminders, err := h.services.Reminder.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	var eventsFull []*models.EventWithFriendsAndReminders

	var eventReminders []models.Reminder
	for _, eventWithFriends := range eventsWithFriends {
		for _, reminder := range reminders {
			if reminder.EventID == eventWithFriends.Event.ID {
				eventReminders = append(eventReminders, reminder)
			}
		}
		eventsFull = append(eventsFull, &models.EventWithFriendsAndReminders{
			Event:     eventWithFriends.Event,
			Friends:   eventWithFriends.Friends,
			Reminders: eventReminders,
		})
		eventReminders = nil
	}

	c.JSON(http.StatusOK, getAllEventsFullInfo{
		Data: eventsFull,
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
