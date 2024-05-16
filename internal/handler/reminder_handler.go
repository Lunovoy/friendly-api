package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lunovoy/friendly/internal/models"
)

// @Summary Create Reminder
// @Security ApiKeyAuth
// @Tags reminder
// @Description create reminder
// @ID create-reminder
// @Accept  json
// @Produce  json
// @Param input body models.Reminder true "Reminder info"
// @Success 201 {string} uuid
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder [post]
func (h *Handler) createReminder(c *gin.Context) {

	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	var payload models.Reminder
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	reminderID, err := h.services.Reminder.Create(userID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"reminder_id": reminderID,
	})

}

// @Summary Get All Reminders
// @Security ApiKeyAuth
// @Tags reminder
// @Description get all reminders
// @ID get-all-reminders
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllRemindersResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder [get]
func (h *Handler) getAllReminders(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	reminders, err := h.services.Reminder.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRemindersResponse{
		Data: reminders,
	})

}

// @Summary Get All Reminders By Event Id
// @Security ApiKeyAuth
// @Tags reminder
// @Description get all reminders by event id
// @ID get-all-reminders-by-event-id
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllRemindersResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder/event/:event_id [get]
func (h *Handler) getAllRemindersByEventID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	eventID, err := uuid.Parse(c.Param("event_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	reminders, err := h.services.Reminder.GetAllByEventID(userID, eventID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRemindersResponse{
		Data: reminders,
	})
}

// @Summary Get Reminder By Id
// @Security ApiKeyAuth
// @Tags reminder
// @Description get reminder by id
// @ID get-reminder-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Reminder
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder/:id [get]
func (h *Handler) getReminderByID(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	reminderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	reminder, err := h.services.Reminder.GetByID(userID, reminderID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"reminder": reminder,
	})
}

// @Summary Update Reminder
// @Security ApiKeyAuth
// @Tags reminder
// @Description update reminder
// @ID update-reminder
// @Accept  json
// @Produce  json
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder/:id [put]
func (h *Handler) updateReminder(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	reminderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	var payload models.ReminderUpdate
	if err := c.BindJSON(&payload); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.services.Reminder.GetByID(userID, reminderID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("reminder not found: %s", err.Error()))
		return
	}

	err = h.services.Reminder.Update(userID, reminderID, payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}

// @Summary Delete Reminder
// @Security ApiKeyAuth
// @Tags reminder
// @Description delete reminder
// @ID delete-reminder
// @Accept  json
// @Produce  json
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/reminder/:id [delete]
func (h *Handler) deleteReminder(c *gin.Context) {
	userID, err := getUserIDFromCtx(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id from ctx not found")
		return
	}

	reminderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	_, err = h.services.Reminder.GetByID(userID, reminderID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("reminder not found or already deleted: %s", err.Error()))
		return
	}

	err = h.services.Reminder.DeleteByID(userID, reminderID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
