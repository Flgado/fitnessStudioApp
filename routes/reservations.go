package routes

import (
	"github.com/Flgado/fitnessStudioApp/handlers"
	"github.com/Flgado/fitnessStudioApp/internal/database/booking"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

func BuildReservationRoutes(dbPoll *sqlx.DB) *chi.Mux {

	// repositories
	readRepo := booking.NewReadRepository(dbPoll)
	wrRepo := booking.NewWriteRepository(dbPoll)

	// usecases
	uc := usecases.NewBookUseCase(readRepo, wrRepo)
	muc := usecases.NewMakeBookUseCase(readRepo, wrRepo)

	h := handlers.NewBookingInfoHandler(uc)
	hm := handlers.NewMakeReservationHandler(muc)

	cRouter := chi.NewRouter()
	cRouter.Get("/users/{userId}/classes", h.HandlerGetUserClasses)
	cRouter.Get("/classes/{classId}/users", h.HandlerGetClassUsers)
	cRouter.Post("/", hm.HandlerCreateBooking)
	return cRouter
}
