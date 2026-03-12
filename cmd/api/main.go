package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"secret-shop/internal/config"
	"secret-shop/internal/database"
	"secret-shop/internal/handler"
	"secret-shop/internal/repository"
	"secret-shop/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.ConnectPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "secret-shop api running"})
	})

	r.GET("/healthz", func(c *gin.Context) {
		pingCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(pingCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "degraded",
				"database": "down",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "up",
		})
	})

	productRepo := repository.NewPostgresProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	customerRepo := repository.NewPostgresCustomerRepository(db)
	authService := service.NewAuthService(customerRepo, cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTokenMinutes)
	authHandler := handler.NewAuthHandler(authService)
	variantRepo := repository.NewPostgresVariantRepository(db)
	variantService := service.NewVariantService(variantRepo)
	variantHandler := handler.NewVariantHandler(variantService)

	apiV1 := r.Group("/api/v1")
	authHandler.RegisterRoutes(apiV1)
	productHandler.RegisterRoutes(apiV1)
	variantHandler.RegisterRoutes(apiV1)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("server listening on port %s", cfg.AppPort)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("start server: %v", err)
	}
}
