package routes

import (
	"github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes(ctx *api.AppContext) chi.Router {
	r := chi.NewRouter()
	authHandler := &handlers.AuthHandler{Ctx: ctx}

	r.Post("/signup", authHandler.HandleSignup)
	r.Post("/login", authHandler.HandleLogin)

	return r
}
