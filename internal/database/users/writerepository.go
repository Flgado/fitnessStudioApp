package users

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/jmoiron/sqlx"
)

type WriteRepository interface {
	Add(ctx context.Context, user api.UpdateUser) error
	Update(ctx context.Context, user api.User) (int64, error)
}

func NewWriteRepository(db *sqlx.DB) WriteRepository {
	return &repository{db: db}
}

func (r *repository) Add(ctx context.Context, user api.UpdateUser) error {
	ur := UserRow{
		Name: user.Name,
	}

	_, err := r.db.NamedExecContext(ctx, AddUserRow, ur)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, user api.User) (int64, error) {
	u := UserRow{}

	row := r.db.QueryRowContext(ctx, findUserById, user.Id)

	err := row.Scan(&u.Id, &u.Name, &u.CreateDate, &u.LastUpdateDate)

	if err != nil {
		return 0, err
	}

	u.Name = user.Name

	rl, err := r.db.NamedExecContext(ctx, UpdateUser, u)
	if err != nil {
		return 0, err
	}

	return rl.RowsAffected()
}
