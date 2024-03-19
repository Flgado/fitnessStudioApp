package booking

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/jmoiron/sqlx"
)

type ReadRepository interface {
	GetUserBookings(ctx context.Context, userId int) ([]api.ClassBooked, error)
	GetClassReservations(ctx context.Context, classId int) ([]api.UsersBooked, error)
	IsClassBookedByUser(ctx context.Context, userId, classId int) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewReadRepository(db *sqlx.DB) ReadRepository {
	return &repository{db: db}
}

func (r *repository) GetUserBookings(ctx context.Context, userId int) ([]api.ClassBooked, error) {

	rows, err := r.db.QueryContext(ctx, GetUserBookings, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bc := []api.ClassBooked{}

	for rows.Next() {
		var classRow ClassBookedRow
		if err = rows.Scan(&classRow.Id, &classRow.Name, &classRow.Date, &classRow.Capacity, &classRow.NumRegistrations, &classRow.ReservedDate); err != nil {
			return nil, err
		}

		// Convert ClassRow to ReadClass
		readClass := api.ClassBooked{
			Id:           classRow.Id,
			Name:         classRow.Name,
			Date:         classRow.Date,
			ReservedDate: classRow.ReservedDate,
		}

		bc = append(bc, readClass)
	}

	if err = rows.Err(); err != nil {
		return nil, nil
	}

	return bc, nil
}

func (r *repository) GetClassReservations(ctx context.Context, classId int) ([]api.UsersBooked, error) {

	rows, err := r.db.QueryContext(ctx, GetUsersOfBooking, classId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ub := []api.UsersBooked{}

	for rows.Next() {
		var userBookRow UserBookedRow
		if err = rows.Scan(&userBookRow.Id, &userBookRow.Name); err != nil {
			return nil, err
		}

		// Convert ClassRow to ReadClass
		readUserBook := api.UsersBooked{
			ClassId:  classId,
			UserId:   userBookRow.Id,
			UserName: userBookRow.Name,
		}

		ub = append(ub, readUserBook)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ub, nil
}

func (r *repository) IsClassBookedByUser(ctx context.Context, userId, classId int) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM booking WHERE user_id = $1 AND class_id = $2", userId, classId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
