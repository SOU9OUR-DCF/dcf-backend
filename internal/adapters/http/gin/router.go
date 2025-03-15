package gin

import (
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/http/gin/handlers"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/http/gin/middleware"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	authService ports.AuthService,
	userService ports.UserService,
	restaurantService ports.RestaurantService,
	eventService ports.EventService,
	volunteerService ports.VolunteerService,
	cfg *config.Config,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(cfg))

	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userService)
	swaggerHandler := handlers.NewSwaggerHandler()
	restaurantHandler := handlers.NewRestaurantHandler(
		restaurantService,
		eventService,
		volunteerService,
	)
	volunteerHandler := handlers.NewVolunteerHandler(volunteerService)
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register_restaurant", authHandler.RegisterRestaurant)
			auth.POST("/register_volunteer", authHandler.RegisterVolunteer)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		users := v1.Group("/user")
		users.Use(authMiddleware.Authenticate())
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/me", userHandler.UpdateMe)
		}

		restaurant := v1.Group("/restaurant")
		restaurant.Use(authMiddleware.Authenticate())
		{
			restaurant.GET("/dashboard", restaurantHandler.GetDashboard)
			restaurant.GET("/", restaurantHandler.GetRestaurant)
			restaurant.POST("/events", restaurantHandler.CreateEvent)
			restaurant.GET("/events/:id", restaurantHandler.GetEvent)
			restaurant.PUT("/events/:id", restaurantHandler.UpdateEvent)
			restaurant.DELETE("/events/:id", restaurantHandler.DeleteEvent)
			restaurant.PATCH("/events/:id/status", restaurantHandler.UpdateEventStatus)
			restaurant.PATCH("/events/:id/guests", restaurantHandler.UpdateGuestCount)
			restaurant.PATCH("/events/:id/meals", restaurantHandler.UpdateMealsServed)

			restaurant.GET("/applications", restaurantHandler.GetVolunteerApplications)
			restaurant.POST("/applications/:id/approve", restaurantHandler.ApproveVolunteerApplication)
			restaurant.POST("/applications/:id/decline", restaurantHandler.DeclineVolunteerApplication)
		}

		volunteer := v1.Group("/volunteer")
		volunteer.Use(authMiddleware.Authenticate())
		{
			volunteer.GET("/dashboard", volunteerHandler.GetVolunteerDashboard)
			volunteer.GET("/upcoming-tasks", volunteerHandler.GetUpcomingTasks)
			volunteer.GET("/nearby-opportunities", volunteerHandler.GetNearbyOpportunities)
			volunteer.GET("/badges", volunteerHandler.GetVolunteerBadges)
			volunteer.POST("/events/:id/apply", volunteerHandler.ApplyForEvent)
			volunteer.POST("/events/:id/check-in", volunteerHandler.CheckInForEvent)
		}
	}

	router.GET("/swagger.yaml", swaggerHandler.SetupSwagger)
	router.GET("/swagger/*any", swaggerHandler.SetupSwaggerUI)

	return router
}
