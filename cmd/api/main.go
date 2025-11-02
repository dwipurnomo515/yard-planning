package main

import (
	"log"
	"net/http"

	"github.com/dwipurnomo515/yard-planning/config"
	"github.com/dwipurnomo515/yard-planning/internal/handler"
	"github.com/dwipurnomo515/yard-planning/internal/middleware"
	"github.com/dwipurnomo515/yard-planning/internal/repository"
	"github.com/dwipurnomo515/yard-planning/internal/service"
	"github.com/dwipurnomo515/yard-planning/pkg/cache"
	"github.com/dwipurnomo515/yard-planning/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewPostgresDB(database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	yardRepo := repository.NewYardRepository(db)
	blockRepo := repository.NewBlockRepository(db)
	planRepo := repository.NewYardPlanRepository(db)
	containerRepo := repository.NewContainerRepository(db)

	// Initialize services
	var containerHandler *handler.ContainerHandler
	var bulkHandler *handler.BulkHandler

	if cfg.EnableCache {
		// Initialize Redis client
		redisClient, err := cache.NewRedisClient(cache.RedisConfig{
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPass,
			DB:       cfg.RedisDB,
		})
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v. Running without cache.", err)
			// Fall back to non-cached service
			containerService := service.NewContainerService(
				yardRepo,
				blockRepo,
				planRepo,
				containerRepo,
			)
			containerHandler = handler.NewContainerHandler(containerService)
			bulkHandler = handler.NewBulkHandler(containerService)
		} else {
			defer redisClient.Close()
			log.Println("Redis cache enabled")

			// Use cached service
			cachedService := service.NewCachedContainerService(
				yardRepo,
				blockRepo,
				planRepo,
				containerRepo,
				redisClient,
			)
			// Convert to base service for handler
			baseService := &cachedService.ContainerService
			containerHandler = handler.NewContainerHandler(baseService)
			bulkHandler = handler.NewBulkHandler(baseService)
		}
	} else {
		log.Println("Cache disabled")
		containerService := service.NewContainerService(
			yardRepo,
			blockRepo,
			planRepo,
			containerRepo,
		)
		containerHandler = handler.NewContainerHandler(containerService)
		bulkHandler = handler.NewBulkHandler(containerService)
	}

	// Setup routes
	mux := http.NewServeMux()

	// Single operation endpoints
	mux.HandleFunc("/suggestion", containerHandler.HandleSuggestion)
	mux.HandleFunc("/placement", containerHandler.HandlePlacement)
	mux.HandleFunc("/pickup", containerHandler.HandlePickup)

	// Bulk operation endpoints (concurrent)
	mux.HandleFunc("/bulk/suggestion", bulkHandler.HandleBulkSuggestion)
	mux.HandleFunc("/bulk/placement", bulkHandler.HandleBulkPlacement)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	handlerWithMiddleware := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(
				middleware.ContentType(mux),
			),
		),
	)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	log.Printf("Cache enabled: %v", cfg.EnableCache)
	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
