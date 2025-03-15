package handlers

import (
	"net/http"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type VolunteerHandler struct {
	volunteerService ports.VolunteerService
}

func NewVolunteerHandler(volunteerService ports.VolunteerService) *VolunteerHandler {
	return &VolunteerHandler{
		volunteerService: volunteerService,
	}
}

func (h *VolunteerHandler) GetVolunteerDashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	dashboard, err := h.volunteerService.GetVolunteerDashboard(c, volunteer.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func (h *VolunteerHandler) GetUpcomingTasks(c *gin.Context) {
	userID := c.GetString("user_id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	tasks, err := h.volunteerService.GetUpcomingTasks(c, volunteer.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *VolunteerHandler) GetNearbyOpportunities(c *gin.Context) {
	userID := c.GetString("user_id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	opportunities, err := h.volunteerService.GetNearbyOpportunities(c, volunteer.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, opportunities)
}

func (h *VolunteerHandler) GetVolunteerBadges(c *gin.Context) {
	userID := c.GetString("user_id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	badges, err := h.volunteerService.GetVolunteerBadges(c, volunteer.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, badges)
}

func (h *VolunteerHandler) ApplyForEvent(c *gin.Context) {
	userID := c.GetString("user_id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	var req struct {
		EventID string `json:"event_id" binding:"required"`
		Role    string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.volunteerService.ApplyForEvent(c, volunteer.ID.String(), req.EventID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

func (h *VolunteerHandler) CheckInForEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	eventVolunteerID := c.Param("id")

	volunteer, err := h.volunteerService.GetVolunteerByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Volunteer not found"})
		return
	}

	err = h.volunteerService.CheckInForEvent(c, volunteer.ID.String(), eventVolunteerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checked in successfully"})
}
