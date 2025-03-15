package application

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
)

type volunteerService struct {
	txManager      ports.TransactionManager
	volunteerRepo  ports.VolunteerRepository
	appRepo        ports.VolunteerApplicationRepository
	eventVolRepo   ports.EventVolunteerRepository
	eventRepo      ports.EventRepository
	restaurantRepo ports.RestaurantRepository
}

func NewVolunteerService(
	txManager ports.TransactionManager,
	volunteerRepo ports.VolunteerRepository,
	appRepo ports.VolunteerApplicationRepository,
	eventVolRepo ports.EventVolunteerRepository,
	eventRepo ports.EventRepository,
	restaurantRepo ports.RestaurantRepository,
) ports.VolunteerService {
	return &volunteerService{
		txManager:      txManager,
		volunteerRepo:  volunteerRepo,
		appRepo:        appRepo,
		eventVolRepo:   eventVolRepo,
		eventRepo:      eventRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *volunteerService) GetVolunteerByUserID(ctx context.Context, userID string) (*domain.Volunteer, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return s.volunteerRepo.GetByUserID(ctx, uid)
}

func (s *volunteerService) GetEventVolunteers(ctx context.Context, eventID string) ([]*domain.Volunteer, error) {
	eid, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	// Get event volunteers
	eventVolunteers, err := s.eventVolRepo.GetByEventID(ctx, eid)
	if err != nil {
		return nil, err
	}

	// Get volunteer details
	volunteers := make([]*domain.Volunteer, 0, len(eventVolunteers))
	for _, ev := range eventVolunteers {
		volunteer, err := s.volunteerRepo.GetByID(ctx, ev.VolunteerID)
		if err != nil {
			return nil, err
		}
		volunteers = append(volunteers, volunteer)
	}

	return volunteers, nil
}

func (s *volunteerService) GetPendingApplications(ctx context.Context, restaurantID string) ([]*domain.VolunteerApplication, error) {
	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	return s.appRepo.GetByRestaurantID(ctx, rid, "pending")
}

func (s *volunteerService) ApproveApplication(ctx context.Context, applicationID string) error {
	appID, err := uuid.Parse(applicationID)
	if err != nil {
		return fmt.Errorf("invalid application ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		// Get the application (read operation, no need for tx)
		app, err := s.appRepo.GetByID(ctx, appID)
		if err != nil {
			return err
		}

		// Update application status with transaction
		if err := s.appRepo.UpdateStatus(ctx, tx, appID, "approved"); err != nil {
			return err
		}

		// Create event volunteer entry with the same transaction
		eventVolunteer := &domain.EventVolunteer{
			EventID:     app.EventID,
			VolunteerID: app.VolunteerID,
			Role:        app.Role,
			CheckedIn:   false,
		}

		return s.eventVolRepo.Create(ctx, tx, eventVolunteer)
	})
}

func (s *volunteerService) DeclineApplication(ctx context.Context, applicationID string) error {
	appID, err := uuid.Parse(applicationID)
	if err != nil {
		return fmt.Errorf("invalid application ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.appRepo.UpdateStatus(ctx, tx, appID, "declined")
	})
}

func (s *volunteerService) GetVolunteerCount(ctx context.Context, restaurantID string) (int, error) {
	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		return 0, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	// Get all events for this restaurant
	events, _, err := s.eventRepo.GetByRestaurantID(ctx, rid, "", 1000, 0)
	if err != nil {
		return 0, err
	}

	// Count unique volunteers across all events
	volunteerMap := make(map[uuid.UUID]bool)
	for _, event := range events {
		eventVolunteers, err := s.eventVolRepo.GetByEventID(ctx, event.ID)
		if err != nil {
			return 0, err
		}

		for _, ev := range eventVolunteers {
			volunteerMap[ev.VolunteerID] = true
		}
	}

	return len(volunteerMap), nil
}

func (s *volunteerService) GetVolunteerDashboard(ctx context.Context, volunteerID string) (map[string]interface{}, error) {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return nil, fmt.Errorf("invalid volunteer ID: %w", err)
	}

	// Get volunteer profile
	volunteer, err := s.volunteerRepo.GetByID(ctx, vid)
	if err != nil {
		return nil, err
	}

	// Get completed tasks count
	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	tasksCompleted := 0
	hoursVolunteered := 0
	mealsServed := 0
	reputationPoints := 0

	for _, ev := range eventVolunteers {
		// Get event details to calculate hours and check if it's past
		event, err := s.eventRepo.GetByID(ctx, ev.EventID)
		if err != nil {
			continue
		}

		if event.Status == domain.EventStatusPast && ev.CheckedIn {
			tasksCompleted++

			// Calculate hours (difference between end and start time in hours)
			duration := event.EndTime.Sub(event.StartTime).Hours()
			hoursVolunteered += int(duration)

			// Add meals served from the event
			mealsServed += event.MealsServed / max(1, len(eventVolunteers)) // Divide by number of volunteers

			// Calculate reputation points (10 per hour)
			reputationPoints += int(duration) * 10
		}
	}

	// Get upcoming tasks
	upcomingTasks, err := s.GetUpcomingTasks(ctx, volunteerID)
	if err != nil {
		return nil, err
	}

	// Get nearby opportunities
	nearbyOpportunities, err := s.GetNearbyOpportunities(ctx, volunteerID)
	if err != nil {
		return nil, err
	}

	// Get earned badges
	badges, err := s.GetVolunteerBadges(ctx, volunteerID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"volunteer":            volunteer,
		"tasks_completed":      tasksCompleted,
		"hours_volunteered":    hoursVolunteered,
		"meals_served":         mealsServed,
		"reputation_points":    reputationPoints,
		"upcoming_tasks":       upcomingTasks,
		"nearby_opportunities": nearbyOpportunities,
		"badges":               badges,
	}, nil
}

func (s *volunteerService) GetUpcomingTasks(ctx context.Context, volunteerID string) ([]map[string]interface{}, error) {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return nil, fmt.Errorf("invalid volunteer ID: %w", err)
	}

	// Get all event volunteers for this volunteer
	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	var upcomingTasks []map[string]interface{}

	for _, ev := range eventVolunteers {
		// Get event details
		event, err := s.eventRepo.GetByID(ctx, ev.EventID)
		if err != nil {
			continue
		}

		// Only include upcoming and active events
		if event.Status == domain.EventStatusUpcoming || event.Status == domain.EventStatusActive {
			// Get restaurant details
			restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
			if err != nil {
				continue
			}

			upcomingTasks = append(upcomingTasks, map[string]interface{}{
				"id":         ev.ID,
				"event_id":   event.ID,
				"title":      event.Title,
				"role":       ev.Role,
				"location":   event.Location,
				"restaurant": restaurant.Name,
				"date":       event.Date,
				"start_time": event.StartTime,
				"end_time":   event.EndTime,
				"status":     event.Status,
				"checked_in": ev.CheckedIn,
				"confirmed":  true, // Assuming if they're in event_volunteers, they're confirmed
			})
		}
	}

	// Get pending applications
	applications, err := s.appRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	for _, app := range applications {
		if app.Status == "pending" {
			// Get event details
			event, err := s.eventRepo.GetByID(ctx, app.EventID)
			if err != nil {
				continue
			}

			// Get restaurant details
			restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
			if err != nil {
				continue
			}

			upcomingTasks = append(upcomingTasks, map[string]interface{}{
				"id":         app.ID,
				"event_id":   event.ID,
				"title":      event.Title,
				"role":       app.Role,
				"location":   event.Location,
				"restaurant": restaurant.Name,
				"date":       event.Date,
				"start_time": event.StartTime,
				"end_time":   event.EndTime,
				"status":     event.Status,
				"checked_in": false,
				"confirmed":  false,
				"pending":    true,
			})
		}
	}

	return upcomingTasks, nil
}

func (s *volunteerService) GetNearbyOpportunities(ctx context.Context, volunteerID string) ([]map[string]interface{}, error) {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return nil, fmt.Errorf("invalid volunteer ID: %w", err)
	}

	// In a real app, we would use geolocation to find nearby events
	// For now, we'll just get upcoming events that the volunteer hasn't applied to yet

	// Get all upcoming events
	// This would need to be implemented in the event repository
	upcomingEvents, err := s.eventRepo.GetUpcomingEvents(ctx)
	if err != nil {
		return nil, err
	}

	// Get all events the volunteer has already applied to or is assigned to
	applications, err := s.appRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	// Create a map of event IDs the volunteer is already involved with
	involvedEvents := make(map[uuid.UUID]bool)
	for _, app := range applications {
		involvedEvents[app.EventID] = true
	}
	for _, ev := range eventVolunteers {
		involvedEvents[ev.EventID] = true
	}

	var nearbyOpportunities []map[string]interface{}

	for _, event := range upcomingEvents {
		// Skip events the volunteer is already involved with
		if involvedEvents[event.ID] {
			continue
		}

		// Get restaurant details
		restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
		if err != nil {
			continue
		}

		// Get current volunteer count
		volunteerCount, err := s.eventVolRepo.CountByEventID(ctx, event.ID)
		if err != nil {
			continue
		}

		// Only include events that still need volunteers
		if volunteerCount < event.MaxVolunteers {
			// Calculate distance (mock for now)
			distance := "1.5 miles away" // In a real app, calculate based on coordinates

			nearbyOpportunities = append(nearbyOpportunities, map[string]interface{}{
				"event_id":          event.ID,
				"title":             event.Title,
				"restaurant_name":   restaurant.Name,
				"location":          event.Location,
				"date":              event.Date,
				"start_time":        event.StartTime,
				"end_time":          event.EndTime,
				"distance":          distance,
				"volunteers_needed": event.MaxVolunteers - volunteerCount,
				"roles_available":   []string{"Food Preparation", "Serving", "Cleanup"},
			})
		}
	}

	return nearbyOpportunities, nil
}

func (s *volunteerService) GetVolunteerBadges(ctx context.Context, volunteerID string) ([]map[string]interface{}, error) {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return nil, fmt.Errorf("invalid volunteer ID: %w", err)
	}

	// Get completed tasks count
	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return nil, err
	}

	// Calculate stats for badge determination
	tasksCompleted := 0
	hoursVolunteered := 0
	mealsServed := 0

	// Track unique roles performed
	roles := make(map[string]bool)

	for _, ev := range eventVolunteers {
		// Get event details
		event, err := s.eventRepo.GetByID(ctx, ev.EventID)
		if err != nil {
			continue
		}

		if event.Status == domain.EventStatusPast && ev.CheckedIn {
			tasksCompleted++

			// Calculate hours
			duration := event.EndTime.Sub(event.StartTime).Hours()
			hoursVolunteered += int(duration)

			// Add meals served
			mealsServed += event.MealsServed / max(1, len(eventVolunteers))

			// Track role
			roles[ev.Role] = true
		}
	}

	// Determine badges based on achievements
	var badges []map[string]interface{}

	// First Timer Badge
	if tasksCompleted >= 1 {
		badges = append(badges, map[string]interface{}{
			"name":        "First Timer",
			"description": "Completed your first volunteer task",
			"earned":      true,
		})
	}

	// Helping Hand Badge
	if tasksCompleted >= 5 {
		badges = append(badges, map[string]interface{}{
			"name":        "Helping Hand",
			"description": "Completed 5 volunteer tasks",
			"earned":      true,
		})
	}

	// Food Server Badge
	if roles["Serving"] {
		badges = append(badges, map[string]interface{}{
			"name":        "Food Server",
			"description": "Served food to those in need",
			"earned":      true,
		})
	}

	// Community Leader Badge
	if tasksCompleted >= 10 {
		badges = append(badges, map[string]interface{}{
			"name":        "Community Leader",
			"description": "Completed 10 volunteer tasks",
			"earned":      true,
		})
	}

	return badges, nil
}

func (s *volunteerService) ApplyForEvent(ctx context.Context, volunteerID string, eventID string, role string) error {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return fmt.Errorf("invalid volunteer ID: %w", err)
	}

	eid, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	// Check if volunteer already applied or is assigned to this event
	applications, err := s.appRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return err
	}

	for _, app := range applications {
		if app.EventID == eid {
			return fmt.Errorf("you have already applied for this event")
		}
	}

	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return err
	}

	for _, ev := range eventVolunteers {
		if ev.EventID == eid {
			return fmt.Errorf("you are already assigned to this event")
		}
	}

	// Get event to check if it's still accepting volunteers
	event, err := s.eventRepo.GetByID(ctx, eid)
	if err != nil {
		return err
	}

	// Check if event is upcoming
	if event.Status != domain.EventStatusUpcoming {
		return fmt.Errorf("this event is not accepting volunteers")
	}

	// Check if event has reached max volunteers
	volunteerCount, err := s.eventVolRepo.CountByEventID(ctx, eid)
	if err != nil {
		return err
	}

	if volunteerCount >= event.MaxVolunteers {
		return fmt.Errorf("this event has reached its volunteer capacity")
	}

	// Create application
	application := &domain.VolunteerApplication{
		ID:          uuid.New(),
		VolunteerID: vid,
		EventID:     eid,
		Role:        role,
		Status:      "pending",
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.appRepo.Create(ctx, tx, application)
	})
}

func (s *volunteerService) CheckInForEvent(ctx context.Context, volunteerID string, eventVolunteerID string) error {
	vid, err := uuid.Parse(volunteerID)
	if err != nil {
		return fmt.Errorf("invalid volunteer ID: %w", err)
	}

	evid, err := uuid.Parse(eventVolunteerID)
	if err != nil {
		return fmt.Errorf("invalid event volunteer ID: %w", err)
	}

	// Verify this event volunteer belongs to this volunteer
	eventVolunteers, err := s.eventVolRepo.GetByVolunteerID(ctx, vid)
	if err != nil {
		return err
	}

	var targetEV *domain.EventVolunteer
	for _, ev := range eventVolunteers {
		if ev.ID == evid {
			targetEV = ev
			break
		}
	}

	if targetEV == nil {
		return fmt.Errorf("event not found for this volunteer")
	}

	// Get event to check if it's active
	event, err := s.eventRepo.GetByID(ctx, targetEV.EventID)
	if err != nil {
		return err
	}

	if event.Status != domain.EventStatusActive {
		return fmt.Errorf("check-in is only available for active events")
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.eventVolRepo.UpdateCheckIn(ctx, tx, evid, true)
	})
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
