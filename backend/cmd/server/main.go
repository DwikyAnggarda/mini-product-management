package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"product-management/backend/internal/config"
	"product-management/backend/internal/db"
	"product-management/backend/internal/handler"
	"product-management/backend/internal/middleware"
	"product-management/backend/internal/repository"
	"product-management/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer func(database *sql.DB) {
		_ = database.Close()
	}(database)

	userRepo := repository.NewUserRepository(database)
	productRepo := repository.NewProductRepository(database)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHour)
	productService := service.NewProductService(productRepo)
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ok"}}`))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/login", authHandler.Login)

		api.Group(func(protected chi.Router) {
			protected.Use(middleware.Auth(cfg.JWTSecret))
			protected.Get("/products", productHandler.List)
			protected.Post("/products", productHandler.Create)
			protected.Put("/products/{id}", productHandler.Update)
			protected.Delete("/products/{id}", productHandler.Delete)
		})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server running on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
