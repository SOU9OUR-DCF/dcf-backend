// cmd/api/main.go
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/http/gin"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/repositories/postgres"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/repositories/redis"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/application"
	"github.com/SOU9OUR-DCF/dcf-backend.git/pkg/jwt"
)

func main() {
	configPath := flag.String("config", "./config", "Path to configuration directory")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbConn, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	redisConn, err := redis.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisConn.Close()

	txManager := postgres.NewTransactionManager(dbConn)

	userRepo := postgres.NewUserRepository(dbConn)
	restaurantRepo := postgres.NewRestaurantRepository(dbConn)
	volunteerRepo := postgres.NewVolunteerRepository(dbConn)
	eventRepo := postgres.NewEventRepository(dbConn)
	volunteerAppRepo := postgres.NewVolunteerApplicationRepository(dbConn)
	eventVolunteerRepo := postgres.NewEventVolunteerRepository(dbConn)
	tokenCache := redis.NewTokenCache(redisConn)

	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	authService := application.NewAuthService(txManager, userRepo, restaurantRepo, volunteerRepo, tokenCache, jwtService)
	userService := application.NewUserService(txManager, userRepo, restaurantRepo, volunteerRepo)
	restaurantService := application.NewRestaurantService(txManager, restaurantRepo, eventRepo, volunteerRepo, volunteerAppRepo, eventVolunteerRepo)
	eventService := application.NewEventService(txManager, eventRepo, restaurantRepo)
	volunteerService := application.NewVolunteerService(txManager, volunteerRepo, volunteerAppRepo, eventVolunteerRepo, eventRepo, restaurantRepo)

	router := gin.NewRouter(
		authService,
		userService,
		restaurantService,
		eventService,
		volunteerService,
		cfg,
	)
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Servers exited properly")
}
