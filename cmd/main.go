package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/api/middleware"
	"github.com/OmerYesilkaya/fileuploader/internal/api/routes"
	"github.com/OmerYesilkaya/fileuploader/internal/config"
	"github.com/OmerYesilkaya/fileuploader/internal/db"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	database, err := db.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Error while initializing db: %v", err)
	}
	defer database.Close()

	ctx := &api.AppContext{
		DB:     database,
		Config: cfg,
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.Cors)

	// Routes
	r.Mount("/auth", routes.AuthRoutes(ctx))
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Mount("/files", routes.FileRoutes(ctx))
	})

	fmt.Printf("Server is running on port %s\n", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Error while setting up server: %v", err)
	}
}
