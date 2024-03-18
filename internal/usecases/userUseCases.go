package usecases

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/database/users"
	"github.com/Flgado/fitnessStudioApp/utils"
)

type UserUseCases interface {
	GetAllUsers(ctx context.Context) ([]api.User, error)
	GetUserById(ctx context.Context, userId int) (api.User, error)
	CreateUser(ctx context.Context, userName string) error
	UpdateUser(ctx context.Context, user api.User) (int64, error)
}

type userUseCases struct {
	readRep  users.ReadRepository
	writeRep users.WriteRepository
}

func NewUserUseCase(readRep users.ReadRepository, writeRepo users.WriteRepository) UserUseCases {
	return &userUseCases{
		readRep:  readRep,
		writeRep: writeRepo,
	}
}

func (u *userUseCases) GetAllUsers(ctx context.Context) ([]api.User, error) {
	return u.readRep.List(ctx)
}

func (u *userUseCases) GetUserById(ctx context.Context, userId int) (api.User, error) {

	user, err := u.readRep.GetById(ctx, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.User{}, utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "User Not Found"},
				"The specified user does not exist.",
				"Please provide a valid class ID.")
		}

		return api.User{}, err
	}

	return user, err
}

func (u *userUseCases) CreateUser(ctx context.Context, userName string) error {
	return u.writeRep.Add(ctx, api.UpdateUser{Name: userName})
}

func (u *userUseCases) UpdateUser(ctx context.Context, user api.User) (int64, error) {
	ur, err := u.writeRep.Update(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "User Not Found"},
				"The specified user does not exist. Unable to update.",
				"Please provide a valid user ID.")
		}

		return 0, err
	}
	return ur, err
}
