package booking

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/jmoiron/sqlx"
)

type WriteRepository interface {
	Add(ctx context.Context, userId int, classId int) error
}

func NewWriteRepository(db *sqlx.DB) WriteRepository {
	return &repository{db: db}
}

func (r *repository) Add(ctx context.Context, userId int, classId int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Lock the row for the specific class being booked
	_, err = tx.ExecContext(ctx, "SELECT * FROM classes WHERE id = $1 FOR UPDATE", classId)
	if err != nil {
		return err
	}

	// Check class capacity
	var numRegistrations, classCapacity int
	err = tx.QueryRowContext(ctx, "SELECT num_registrations, class_capacity FROM classes WHERE id = $1", classId).Scan(&numRegistrations, &classCapacity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "Class Not Found"},
				"The specified class does not exist.",
				"Please provide a valid class ID.")
		}

		return err
	}

	if numRegistrations >= classCapacity {
		return utils.E(http.StatusUnprocessableEntity,
			nil,
			map[string]string{"message": "Class Capacity Reached"},
			"The class is already full and cannot accept any more registrations.",
			"Please select another class or try again later.")
	}

	// Increment num_registrations
	_, err = tx.ExecContext(ctx, "UPDATE classes SET num_registrations = num_registrations + 1 WHERE id = $1", classId)
	if err != nil {
		return err
	}

	// Insert booking record
	_, err = tx.ExecContext(ctx, "INSERT INTO booking (user_id, class_id, reserved_date) VALUES ($1, $2, CURRENT_TIMESTAMP)", userId, classId)
	if err != nil {
		return err
	}

	return nil
}
