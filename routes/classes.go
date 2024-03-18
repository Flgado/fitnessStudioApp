package routes

import (
	"github.com/Flgado/fitnessStudioApp/handlers"
	"github.com/Flgado/fitnessStudioApp/internal/database/classes"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func BuildClassesRoutes(dbPoll *sqlx.DB) *chi.Mux {

	// repositories
	readRepo := classes.NewReadRepository(dbPoll)
	wrRepo := classes.NewWriteRepository(dbPoll)

	// usecases
	uc := usecases.NewClassesUseCases(readRepo, wrRepo)

	// handler
	h := handlers.NewClassesHandler(uc)
	cRouter := chi.NewRouter()
	cRouter.Get("/", h.HandlerGetClasses)
	cRouter.Get("/{class-id}", h.HandlerGetClassById)
	cRouter.Post("/", h.HandlerAddClass)
	cRouter.Post("/{class-id}", h.HandlerUpdateClass)

	return cRouter
}
