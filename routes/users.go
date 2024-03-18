package routes

import (
	"github.com/Flgado/fitnessStudioApp/handlers"
	"github.com/Flgado/fitnessStudioApp/internal/database/users"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func BuildUserRoutes(dbPoll *sqlx.DB) *chi.Mux {
	// repository
	rr := users.NewReadRepository(dbPoll)
	wr := users.NewWriteRepository(dbPoll)

	// usecases
	gu := usecases.NewUserUseCase(rr, wr)

	// handlers
	h := handlers.NewUsersHandler(gu)
	uRouter := chi.NewRouter()
	uRouter.Get("/", h.HandlerGetUsers)
	uRouter.Get("/{user-id}", h.HandlerGetUserById)
	uRouter.Post("/", h.HandlerCreateUser)
	uRouter.Post("/{user-id}", h.HandlerUpdateUser)
	return uRouter
}
