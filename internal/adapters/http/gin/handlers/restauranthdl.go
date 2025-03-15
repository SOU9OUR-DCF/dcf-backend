package handlers

import (
	"net/http"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	restaurantService ports.RestaurantService
	eventService      ports.EventService
	volunteerService  ports.VolunteerService
}

func NewRestaurantHandler(
	restaurantService ports.RestaurantService,
	eventService ports.EventService,
	volunteerService ports.VolunteerService,
) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		eventService:      eventService,
		volunteerService:  volunteerService,
	}
}

func (h *RestaurantHandler) GetDashboard(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant stats
	stats, err := h.restaurantService.GetRestaurantStats(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get upcoming events
	upcomingEvents, _, err := h.eventService.GetUpcomingEvents(c.Request.Context(), restaurant.ID.String(), 5, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get today's events
	todayEvents, err := h.eventService.GetTodayEvents(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pending volunteer applications
	pendingApps, err := h.volunteerService.GetPendingApplications(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant":      restaurant,
		"stats":           stats,
		"upcoming_events": upcomingEvents,
		"today_events":    todayEvents,
		"pending_apps":    pendingApps,
	})
}

func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get restaurant stats
	stats, err := h.restaurantService.GetRestaurantStats(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get upcoming events
	upcomingEvents, _, err := h.eventService.GetUpcomingEvents(c.Request.Context(), restaurant.ID.String(), 5, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get today's events
	todayEvents, err := h.eventService.GetTodayEvents(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pending volunteer applications
	pendingApps, err := h.volunteerService.GetPendingApplications(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant":      restaurant,
		"stats":           stats,
		"upcoming_events": upcomingEvents,
		"today_events":    todayEvents,
		"pending_apps":    pendingApps,
	})
}

func (h *RestaurantHandler) CreateEvent(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var event domain.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set restaurant ID
	event.RestaurantID = restaurant.ID

	if err := h.eventService.CreateEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *RestaurantHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get volunteers for this event
	volunteers, err := h.volunteerService.GetEventVolunteers(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event":      event,
		"volunteers": volunteers,
	})
}

func (h *RestaurantHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")

	// Get the existing event
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if event.RestaurantID != restaurant.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this event"})
		return
	}

	// Bind updated event data
	var updatedEvent domain.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve the ID and restaurant ID
	updatedEvent.ID = event.ID
	updatedEvent.RestaurantID = event.RestaurantID

	if err := h.eventService.UpdateEvent(c.Request.Context(), &updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

func (h *RestaurantHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	// Get the existing event
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if event.RestaurantID != restaurant.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to delete this event"})
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}

func (h *RestaurantHandler) UpdateEventStatus(c *gin.Context) {
	eventID := c.Param("id")

	// Get the existing event
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if event.RestaurantID != restaurant.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this event"})
		return
	}

	var req struct {
		Status domain.EventStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.eventService.UpdateEventStatus(c.Request.Context(), eventID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event status updated successfully"})
}

func (h *RestaurantHandler) UpdateGuestCount(c *gin.Context) {
	eventID := c.Param("id")

	// Get the existing event
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if event.RestaurantID != restaurant.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this event"})
		return
	}

	var req struct {
		Count int `json:"count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.eventService.UpdateGuestCount(c.Request.Context(), eventID, req.Count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guest count updated successfully"})
}

func (h *RestaurantHandler) UpdateMealsServed(c *gin.Context) {
	eventID := c.Param("id")

	// Get the existing event
	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	if event.RestaurantID != restaurant.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this event"})
		return
	}

	var req struct {
		Count int `json:"count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.eventService.UpdateMealsServed(c.Request.Context(), eventID, req.Count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meals served count updated successfully"})
}

func (h *RestaurantHandler) GetVolunteerApplications(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get restaurant by user ID
	restaurant, err := h.restaurantService.GetRestaurantByUserID(c.Request.Context(), user.(*domain.User).ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applications, err := h.volunteerService.GetPendingApplications(c.Request.Context(), restaurant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}

func (h *RestaurantHandler) ApproveVolunteerApplication(c *gin.Context) {
	applicationID := c.Param("id")

	if err := h.volunteerService.ApproveApplication(c.Request.Context(), applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application approved successfully"})
}

func (h *RestaurantHandler) DeclineVolunteerApplication(c *gin.Context) {
	applicationID := c.Param("id")

	if err := h.volunteerService.DeclineApplication(c.Request.Context(), applicationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application declined successfully"})
}
