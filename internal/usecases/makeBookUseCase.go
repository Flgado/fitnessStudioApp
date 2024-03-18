package usecases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Flgado/fitnessStudioApp/internal/database/booking"
	"github.com/Flgado/fitnessStudioApp/utils"
)

type MakeBookUseCase interface {
	Book(ctx context.Context, userId int, classId int) error
}

type makeBookUseCase struct {
	readRep booking.ReadRepository
	wrRep   booking.WriteRepository
}

func NewMakeBookUseCase(readRep booking.ReadRepository, wrRep booking.WriteRepository) MakeBookUseCase {
	return &makeBookUseCase{
		readRep: readRep,
		wrRep:   wrRep,
	}
}

func (uc *makeBookUseCase) Book(ctx context.Context, userId int, classId int) error {
	reserved, err := uc.readRep.IsClassBookedByUser(ctx, userId, classId)

	if err != nil {
		return err
	}

	if reserved {
		return utils.E(http.StatusConflict,
			nil,
			map[string]string{"message": "Conflict Status"},
			fmt.Sprintf("Class with Id: %d is already reserved by User with id %d",
				userId,
				classId),
			"Validate user reserved classes")
	}

	err = uc.wrRep.Add(ctx, userId, classId)

	if err != nil {
		return err
	}

	return nil
}
