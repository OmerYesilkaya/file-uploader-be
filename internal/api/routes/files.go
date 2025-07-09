package routes

import (
	// "github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func FileRoutes(ctx *api.AppContext) chi.Router {
	r := chi.NewRouter()
	fileHandler := &handlers.FileHandler{Ctx: ctx}

	r.Post("/upload", fileHandler.HandleFileUpload)

	return r
}
