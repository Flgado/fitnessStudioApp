package users

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	ReadRepository
	WriteRepository
}

type ReadRepository interface {
	List(ctx context.Context) ([]api.User, error)
	GetById(ctx context.Context, id int) (api.User, error)
	GetByName(ctx context.Context, name string) ([]api.User, error)
}

type repository struct {
	db *sqlx.DB
}

func NewReadRepository(db *sqlx.DB) ReadRepository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]api.User, error) {
	rows, err := r.db.QueryxContext(ctx, findUsers)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.QueryxContext")
	}

	defer rows.Close()

	u := []api.User{}

	for rows.Next() {
		var user UserRow
		if err = rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "authRepo.FindByName.StructScan")
		}

		readUser := api.User{
			Id:   user.Id,
			Name: user.Name,
		}
		u = append(u, readUser)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.rows.Err")
	}

	return u, nil
}

func (r *repository) GetById(ctx context.Context, userId int) (api.User, error) {
	u := UserRow{}

	row := r.db.QueryRowContext(ctx, findUserById, userId)
	err := row.Scan(&u.Id, &u.Name, &u.CreateDate, &u.LastUpdateDate)

	if err != nil {
		return api.User{}, err
	}

	readUser := api.User{
		Id:   u.Id,
		Name: u.Name,
	}

	return readUser, nil
}

func (r *repository) GetByName(ctx context.Context, userName string) ([]api.User, error) {
	rows, err := r.db.QueryxContext(ctx, findUserByName, userName)

	if err != nil {
		return nil, err
	}

	u := []api.User{}

	for rows.Next() {
		var user api.User
		if err = rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "authRepo.FindByName.StructScan")
		}
		u = append(u, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.rows.Err")
	}

	return u, nil
}
