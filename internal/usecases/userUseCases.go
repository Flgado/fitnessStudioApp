package usecases

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/database/users"
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

func NewGetAllUserUseCase(readRep users.ReadRepository, writeRepo users.WriteRepository) UserUseCases {
	return &userUseCases{
		readRep:  readRep,
		writeRep: writeRepo,
	}
}

func (u *userUseCases) GetAllUsers(ctx context.Context) ([]api.User, error) {
	return u.readRep.List(ctx)
}

func (u *userUseCases) GetUserById(ctx context.Context, userId int) (api.User, error) {
	return u.readRep.GetById(ctx, userId)
}

func (u *userUseCases) CreateUser(ctx context.Context, userName string) error {
	return u.writeRep.Add(ctx, api.UpdateUser{Name: userName})
}

func (u *userUseCases) UpdateUser(ctx context.Context, user api.User) (int64, error) {
	return u.writeRep.Update(ctx, user)
}
