package classes

import (
	"context"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/jmoiron/sqlx"
)

type ReadRepository interface {
	List(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error)
	GetById(ctx context.Context, classId int) (api.ReadClass, error)
	GetClassReservations(ctx context.Context, classId int) (int, error)
}

type repository struct {
	db *sqlx.DB
}

func NewReadRepository(db *sqlx.DB) ReadRepository {
	return &repository{db: db}
}

// List retrieves classes from the repository based on the provided filters.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and a api.ClasseFilters struct containing optional filtering parameters for classes.
// It constructs a SQL query based on the provided filters and executes it against the database.
// It returns a slice of api.ReadClass structs representing the retrieved classes and nil error if successful.
// If there is an issue retrieving the classes from the database, it returns nil slice and an error describing the issue.
//
// param: ctx context.Context - Context object for managing the request lifecycle.
// param: filters api.ClasseFilters - Struct containing optional filtering parameters for classes.
//
// @return []api.ReadClass - Slice of ReadClass structs representing the retrieved classes.
// @return error - Error if there is an issue retrieving the classes from the database.
func (r *repository) List(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error) {
	query := "SELECT * FROM classes WHERE 1=1"

	args := make(map[string]interface{})
	if filters.Name != "" {
		query += " AND class_name = :name"
		args["name"] = filters.Name
	}
	if filters.StartDateGte != nil {
		query += " AND class_date >= :start_date_gte"
		args["start_date_gte"] = *filters.StartDateGte
	}
	if filters.EndDateLe != nil {
		query += " AND class_date <= :end_date_le"
		args["end_date_le"] = *filters.EndDateLe
	}
	if filters.CapacityGte != nil {
		query += " AND class_capacity >= :capacity_gte"
		args["capacity_gte"] = *filters.CapacityGte
	}
	if filters.CapacityLe != nil {
		query += " AND class_capacity <= :capacity_le"
		args["capacity_le"] = *filters.CapacityLe
	}
	if filters.NumRegistrationsGte != nil {
		query += " AND num_registrations >= :num_registrations_gte"
		args["num_registrations_gte"] = filters.NumRegistrationsGte
	}
	if filters.NumRegistrationsLe != nil {
		query += " AND num_registrations <= :num_registrations_le"
		args["num_registrations_le"] = filters.NumRegistrationsLe
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	classes := []api.ReadClass{}

	for rows.Next() {
		var classRow ClassRow
		if err = rows.StructScan(&classRow); err != nil {
			return nil, err
		}
		// Convert ClassRow to ReadClass
		readClass := api.ReadClass{
			Id: classRow.Id,
			Class: api.Class{
				Name:     classRow.Name,
				Date:     classRow.Date,
				Capacity: classRow.Capacity,
			},
			NumRegistrations: classRow.NumRegistrations,
		}
		classes = append(classes, readClass)
	}

	if err = rows.Err(); err != nil {
		return nil, nil
	}

	return classes, nil
}

func (r *repository) GetById(ctx context.Context, classId int) (api.ReadClass, error) {
	cr := ClassRow{}
	row := r.db.QueryRowContext(ctx, findClassById, classId)

	err := row.Scan(&cr.Id, &cr.Name, &cr.Date, &cr.Capacity, &cr.NumRegistrations, &cr.CreateDate, &cr.LastUpdateDate)

	if err != nil {
		return api.ReadClass{}, err
	}

	readClass := api.ReadClass{
		Id: cr.Id,
		Class: api.Class{
			Name:     cr.Name,
			Date:     cr.Date,
			Capacity: cr.Capacity,
		},
		NumRegistrations: cr.NumRegistrations,
	}

	return readClass, nil
}

func (r *repository) GetClassReservations(ctx context.Context, classId int) (int, error) {
	var c int

	row := r.db.QueryRowContext(ctx, classReservationsById, classId)

	err := row.Scan(c)

	if err != nil {
		return 0, err
	}

	return c, nil
}
