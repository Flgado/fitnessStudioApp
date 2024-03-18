package usecases

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/database/booking"
)

type BookingUseCase interface {
	GetUserReservations(ctx context.Context, userId int) ([]api.ClassBooked, error)
	GetClassesReservations(ctx context.Context, classId int) ([]api.UsersBooked, error)
}

type bookingUseCases struct {
	readRep booking.ReadRepository
	wrRep   booking.WriteRepository
}

func NewBookUseCase(readRepo booking.ReadRepository, wrRepo booking.WriteRepository) BookingUseCase {
	return &bookingUseCases{
		readRep: readRepo,
		wrRep:   wrRepo,
	}
}

func (uc *bookingUseCases) GetUserReservations(ctx context.Context, userId int) ([]api.ClassBooked, error) {
	return uc.readRep.GetUserBookings(ctx, userId)
}

func (uc *bookingUseCases) GetClassesReservations(ctx context.Context, classId int) ([]api.UsersBooked, error) {
	return uc.readRep.GetClassReservations(ctx, classId)
}
