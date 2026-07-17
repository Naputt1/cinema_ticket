package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"github.com/naputt/cinema-ticket-backend/internal/events"
	"github.com/naputt/cinema-ticket-backend/internal/handler"
	mongoRepo "github.com/naputt/cinema-ticket-backend/internal/repository/mongo"
	redisRepo "github.com/naputt/cinema-ticket-backend/internal/repository/redis"
	"github.com/naputt/cinema-ticket-backend/internal/router"
	"github.com/naputt/cinema-ticket-backend/internal/usecase"
	"github.com/naputt/cinema-ticket-backend/pkg/config"
)

// @title Cinema Ticket Booking API
// @version 1.0
// @description Cinema ticket booking system with real-time seat updates via SSE.
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(cfg.Mongo.DB)

	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	rmqPublisher, err := events.NewRabbitMQPublisher(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rmqPublisher.Close()

	userRepo := mongoRepo.NewUserRepository(mongoDB)
	movieRepo := mongoRepo.NewMovieRepository(mongoDB)
	showtimeRepo := mongoRepo.NewShowtimeRepository(mongoDB)
	seatRepo := mongoRepo.NewSeatRepository(mongoDB)
	bookingRepo := mongoRepo.NewBookingRepository(mongoDB)
	auditRepo := mongoRepo.NewAuditLogRepository(mongoDB)
	lockRepo := redisRepo.NewLockRepository(redisClient)

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	movieUsecase := usecase.NewMovieUsecase(movieRepo, showtimeRepo, seatRepo)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, seatRepo, showtimeRepo, lockRepo, rmqPublisher, auditRepo)
	adminUsecase := usecase.NewAdminUsecase(bookingRepo, auditRepo)
	sseHandler := handler.NewSSEHandler()

	// RabbitMQ consumer for SSE bridge
	sseConsumer, err := events.NewRabbitMQConsumer(
		cfg.RabbitMQ.URL,
		"seat.events",
		[]string{"seat.*"},
		func(event domain.Event) error {
			sseHandler.HandleEvent(event)
			return nil
		},
	)
	if err != nil {
		log.Printf("warning: failed to create SSE consumer: %v", err)
	} else {
		if err := sseConsumer.Start(); err != nil {
			log.Printf("warning: failed to start SSE consumer: %v", err)
		}
		defer sseConsumer.Stop()
	}

	// Notification consumer (mock)
	notifConsumer, err := events.NewRabbitMQConsumer(
		cfg.RabbitMQ.URL,
		"notification",
		[]string{"booking.*", "system.error"},
		func(event domain.Event) error {
			log.Printf("[MOCK NOTIFICATION] event=%s payload=%+v", event.Type, event.Payload)
			return nil
		},
	)
	if err != nil {
		log.Printf("warning: failed to create notification consumer: %v", err)
	} else {
		if err := notifConsumer.Start(); err != nil {
			log.Printf("warning: failed to start notification consumer: %v", err)
		}
		defer notifConsumer.Stop()
	}

	// Booking expiry cleanup
	cleanup := usecase.NewBookingExpiryCleanup(bookingRepo, seatRepo, lockRepo, rmqPublisher, auditRepo)
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	go cleanup.Start(cleanupCtx, 30*time.Second)

	r := gin.Default()
	router.Setup(r, &router.Handlers{
		Auth:    handler.NewAuthHandler(authUsecase, cfg.App.FrontendURL),
		Movie:   handler.NewMovieHandler(movieUsecase),
		Booking: handler.NewBookingHandler(bookingUsecase),
		Admin:   handler.NewAdminHandler(adminUsecase),
		SSE:     sseHandler,
	}, cfg.JWT.Secret, cfg.App.FrontendURL)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cleanupCancel()
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
