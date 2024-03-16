package classes

import (
	"context"
	"strings"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/jmoiron/sqlx"
)

type WriteRepository interface {
	Add(ctx context.Context, user []api.Class) error
	Update(ctx context.Context, classId int, classUpdate api.UpdateClass) (int64, error)
	//AddMultiple(ctx context.Context, users []api.User) error
}

func NewWriteRepository(db *sqlx.DB) WriteRepository {
	return &repository{db: db}
}

// Add inserts classes into the repository.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and a slice of api.Class structs representing the classes to be inserted.
// It converts each api.Class to a ClassRow struct and executes a bulk insert operation
// into the database using the NamedExecContext method.
// It returns nil if the classes are successfully inserted, otherwise, it returns an error.
//
// @param ctx context.Context - Context object for managing the request lifecycle.
// @param classes []api.Class - Slice of api.Class structs representing the classes to be inserted.
//
// @return error - Error if there is an issue inserting the classes into the database.
func (r *repository) Add(ctx context.Context, classes []api.Class) error {
	// Prepare slice to hold ClassRow values
	classRows := make([]ClassRow, len(classes))

	// Convert each api.Class to ClassRow
	for i, class := range classes {
		classRows[i] = ClassRow{
			Name:     class.Name,
			Date:     class.Date,
			Capacity: class.Capacity,
		}
	}

	// Execute bulk insert
	_, err := r.db.NamedExecContext(ctx, AddClassRow, classRows)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing class in the repository.
//
// This method takes a context.Context object for managing the lifecycle of the request
// an integer representing the ID of the class to be updated, and an api.UpdateClass struct
// containing the fields to be modified.
// It fetches the existing class data from the database, constructs a SQL UPDATE query based
// on the provided update fields, and executes it using the NamedExecContext method.
// It returns the number of rows affected if the update is successful, otherwise, it returns an error.
//
// @param ctx context.Context - Context object for managing the request lifecycle.
// @param classId int - ID of the class to be updated.
// @param classUpdate api.UpdateClass - Struct containing the fields to be modified.
//
// @return int64 - Number of rows affected by the update operation.
// @return error - Error if there is an issue updating the class in the database.
func (r *repository) Update(ctx context.Context, classId int, classUpdate api.UpdateClass) (int64, error) {
	// Fetch existing class data from the database
	existingClass := ClassRow{}
	err := r.db.GetContext(ctx, &existingClass, findClassById, classId)

	// if class does not exist
	if err != nil {
		return 0, err
	}

	query := "UPDATE classes SET "
	args := map[string]interface{}{
		"id": classId,
	}
	var updateFields []string

	if classUpdate.Name != nil {
		updateFields = append(updateFields, "class_name=:name")
		args["name"] = *classUpdate.Name
	}
	if classUpdate.Date != nil {
		updateFields = append(updateFields, "start_date=:start_date")
		args["start_date"] = *classUpdate.Date
	}

	if classUpdate.Capacity != nil {
		updateFields = append(updateFields, "class_capacity=:capacity")
		args["capacity"] = *classUpdate.Capacity
	}

	// Check if any fields are to be updated
	if len(updateFields) == 0 {
		return 0, nil // No fields to update
	}

	// Append all update fields to the query
	query += strings.Join(updateFields, ", ")
	query += " WHERE id=:id"

	// Execute the update query
	result, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return 0, err
	}

	// Return the number of rows affected
	return result.RowsAffected()
}
