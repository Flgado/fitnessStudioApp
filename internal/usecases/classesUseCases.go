package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/database/classes"
	"github.com/Flgado/fitnessStudioApp/utils"
)

type reservedDaysInfo struct {
	days []time.Time
	mu   sync.Mutex
}

type ClassesUseCases interface {
	GetFilteredClasses(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error)
	CreateClass(ctx context.Context, class api.ClassScheduler) ([]api.Class, error)
	UpdateClass(ctx context.Context, updateClass api.UpdateClass, classId int) (int64, error)
	GetClassById(ctx context.Context, classId int) (api.ReadClass, error)
}

type classesUseCases struct {
	readRep      classes.ReadRepository
	wrRep        classes.WriteRepository
	reservedDays sync.Map
}

func NewClassesUseCases(readRepo classes.ReadRepository, wrRepo classes.WriteRepository) ClassesUseCases {
	return &classesUseCases{
		readRep: readRepo,
		wrRep:   wrRepo,
	}
}

// GetFilteredClasses retrieves a list of classes filtered by the provided filters.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and a api.ClasseFilters struct containing optional filtering parameters for classes.
// It returns a slice of api.ReadClass structs representing the filtered classes and nil error if successful.
// If there is an issue retrieving the filtered classes, it returns an empty slice and an error describing the issue.
//
// param: ctx context.Context - Context object for managing the request lifecycle.
// param: filters api.ClasseFilters - Struct containing optional filtering parameters for classes.
//
// @return []api.ReadClass - Slice of ReadClass structs representing the filtered classes.
// @return error - Error if there is an issue retrieving the filtered classes.
func (c *classesUseCases) GetFilteredClasses(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error) {
	return c.readRep.List(ctx, filters)
}

// GetClassById retrieves a class by its unique identifier.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and an integer representing the ID of the class to retrieve.
// It returns a api.ReadClass struct representing the class if found and nil error.
// If the class with the specified ID does not exist, it returns an error with a HTTP 404 status code.
// If there is an issue retrieving the class, it returns an empty api.ReadClass struct and an error describing the issue.
//
// param: ctx context.Context - Context object for managing the request lifecycle.
// param: classId int - ID of the class to retrieve.
//
// @return api.ReadClass - ReadClass struct representing the class.
// @return error - Error if there is an issue retrieving the class.
func (c *classesUseCases) GetClassById(ctx context.Context, classId int) (api.ReadClass, error) {
	class, err := c.readRep.GetById(ctx, classId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.ReadClass{}, utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "Class Not Found"},
				"The specified class does not exist. Unable to update.",
				"Please provide a valid class ID.")
		}
		return api.ReadClass{}, err
	}

	return class, nil
}

// CreateClass creates classes based on the provided class scheduler and adds them to the repository.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and a api.ClassScheduler struct containing details about the classes to be created.
// It separates the classes by year and month, checks for availability, and adds them to the repository.
// It returns a slice of api.Class structs representing the classes that could not be scheduled
// due to unavailability or errors, and nil error if successful.
//
// param: ctx context.Context - Context object for managing the request lifecycle.
// param: classScheduler api.ClassScheduler - Struct containing details about the classes to be created.
//
// @return []api.Class - Slice of Class structs representing the classes that could not be scheduled.
// @return error - Error if there is an issue scheduling the classes.
func (c *classesUseCases) CreateClass(ctx context.Context, classScheduler api.ClassScheduler) ([]api.Class, error) {

	if classScheduler.EndDate.Before(classScheduler.StartDate) {
		return []api.Class{}, utils.E(http.StatusBadRequest,
			nil,
			map[string]string{"message": "BadRequest"},
			"End Date should be higher or equals then Start Date",
			"Please select the dates accurately.")

	}
	sc := separateClassByYearMonth(classScheduler)
	var notPossibleSchedulerReport []api.Class
	for key, classList := range sc {

		possibleScheduler, impossibleToSheduler, err := c.getAvailableDays(key, classList)

		if len(impossibleToSheduler) != 0 {
			notPossibleSchedulerReport = append(notPossibleSchedulerReport, impossibleToSheduler...)
		}
		if err != nil {
			// all classes cannot be scheduler
			return append(possibleScheduler, notPossibleSchedulerReport...), err
		}

		if len(possibleScheduler) != 0 {
			err = c.wrRep.Add(ctx, possibleScheduler)
		}

		if err != nil {
			// Remove the values from the cache if something went wrong in the repository
			_ = c.removeDaysFromCache(key, possibleScheduler)
			// all classes cannot be scheduler
			return append(possibleScheduler, notPossibleSchedulerReport...), err
		}
	}

	return notPossibleSchedulerReport, nil
}

// UpdateClass updates the details of a class with the provided information.
//
// This method takes a context.Context object for managing the lifecycle of the request
// and an api.UpdateClass struct containing the updated details of the class.
// It also takes the ID of the class to be updated.
// It performs validations such as checking if the provided date is in the past
// and if the updated capacity can be set, and then updates the class in the repository.
// It returns the number of rows affected by the update operation and nil error if successful.
//
// param: ctx context.Context - Context object for managing the request lifecycle.
// param: updateClass api.UpdateClass - Struct containing the updated details of the class.
// param: classId int - ID of the class to be updated.
//
// @return int64 - Number of rows affected by the update operation.
// @return error - Error if there is an issue updating the class.
func (c *classesUseCases) UpdateClass(ctx context.Context, updateClass api.UpdateClass, classId int) (int64, error) {

	if updateClass.Date != nil && updateClass.Date.Before(time.Now()) {
		return 0, utils.E(http.StatusUnprocessableEntity,
			nil,
			map[string]string{"message": "Status Unprocessabe Entity"},
			"New Date cannot be in the pass",
			"Please select a valid day")
	}

	_, err := c.readRep.GetById(ctx, classId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "Class Not Found"},
				"The specified class does not exist. Unable to update.",
				"Please provide a valid class ID.")
		}

		return 0, err
	}

	if updateClass.Date != nil {
		// Validate in cache if day is availabe
		key := fmt.Sprintf("%d-%02d", updateClass.Date.Year(), updateClass.Date.Month())

		isAvailable := c.isDayAvailable(key, *updateClass.Date)
		if !isAvailable {
			return 0, utils.E(http.StatusNotFound,
				nil,
				map[string]string{"message": "Date already reserved"},
				"The selected date is already reserved.",
				"Please choose a different date or class.")
		}
	}

	return c.wrRep.Update(ctx, classId, updateClass)
}

// removeDaysFromCache removes reserved days from the cache for a specific month.
//
// This method takes a string key representing the month and a slice of api.Class
// representing the classes whose reserved days need to be removed from the cache.
// It removes the reserved days associated with the provided classes from the cache.
// It returns nil if the operation is successful.
//
// param: key string - Key representing the month (e.g., "2024-03").
// param: classList []api.Class - Slice of Class structs representing the classes.
//
// @return error - Error if there is an issue removing reserved days from the cache.
func (c *classesUseCases) removeDaysFromCache(key string, classList []api.Class) error {
	value, _ := c.reservedDays.LoadOrStore(key, &reservedDaysInfo{})
	info := value.(*reservedDaysInfo)

	// Lock to prevent concurrent access to the reserved days slice
	info.mu.Lock()
	defer info.mu.Unlock()

	// Create a map of reserved days for constant-time lookup
	daysToRemove := make(map[int]struct{})
	for _, day := range classList {
		daysToRemove[day.Date.Day()] = struct{}{}
	}

	var newMouthCache []time.Time
	for _, cached := range info.days {
		if _, reserved := daysToRemove[cached.Day()]; !reserved {
			newMouthCache = append(newMouthCache, cached)
		}
	}

	c.reservedDays.Store(key, newMouthCache)
	return nil
}

// isDayAvailable checks if a specific day is available for scheduling.
//
// This method takes a string key representing the month and a time.Time object
// representing the day to be checked for availability.
// It checks if the provided day is already reserved and returns true if it is available,
// otherwise returns false.
//
// param: key string - Key representing the month (e.g., "2024-03").
// param: date time.Time - Date to be checked for availability.
//
// @return bool - True if the day is available, false otherwise.
func (c *classesUseCases) isDayAvailable(key string, date time.Time) bool {
	// Load or initialize reserved days info for the key
	value, _ := c.reservedDays.LoadOrStore(key, &reservedDaysInfo{})
	info := value.(*reservedDaysInfo)

	// Lock to prevent concurrent access to the reserved days slice
	info.mu.Lock()
	defer info.mu.Unlock()

	// Create a map of reserved days for constant-time lookup
	reservedMap := make(map[int]struct{})
	for _, day := range info.days {
		reservedMap[day.Day()] = struct{}{}
	}

	if _, reserved := reservedMap[date.Day()]; reserved {
		return false
	}

	info.days = append(info.days, date)

	return true
}

// getAvailableDays retrieves available days for scheduling classes based on the provided class list.
//
// This method takes a string key representing the month and a slice of api.Class
// representing the classes to be scheduled.
// It checks for available days within the month and returns a slice of available classes
// and a slice of classes that could not be scheduled due to unavailability.
// It returns nil error if successful.
//
// param: key string - Key representing the month (e.g., "2024-03").
// param: classList []api.Class - Slice of Class structs representing the classes to be scheduled.
//
// @return []api.Class - Slice of available Class structs.
// @return []api.Class - Slice of unavailable Class structs.
// @return error - Error if there is an issue retrieving available days.
func (c *classesUseCases) getAvailableDays(key string, classList []api.Class) ([]api.Class, []api.Class, error) {
	// Load or initialize reserved days info for the key
	value, _ := c.reservedDays.LoadOrStore(key, &reservedDaysInfo{})
	info := value.(*reservedDaysInfo)

	// Lock to prevent concurrent access to the reserved days slice
	info.mu.Lock()
	defer info.mu.Unlock()

	reservedMap := make(map[int]struct{})
	for _, day := range info.days {
		reservedMap[day.Day()] = struct{}{}
	}

	// Filter out the reserved days from the classList
	var availableDays []api.Class
	var notPossibleToReserve []api.Class

	for _, class := range classList {
		if _, reserved := reservedMap[class.Date.Day()]; !reserved {
			availableDays = append(availableDays, class)
		} else {
			notPossibleToReserve = append(notPossibleToReserve, class)
		}
	}

	// Update reserved days in the sync.Map
	for _, class := range availableDays {
		info.days = append(info.days, class.Date)
	}

	c.reservedDays.Store(key, info)

	// Return the available days for reservation
	return availableDays, notPossibleToReserve, nil
}

// separateClassByYearMonth separates classes by year and month based on the provided class scheduler.
//
// This method takes an api.ClassScheduler struct representing the classes to be scheduled
// and separates them into a map where the keys are strings representing the year and month (e.g., "2024-03")
// and the values are slices of api.Class representing the classes scheduled for each month.
// It returns the map containing the separated classes.
//
// param: base api.ClassScheduler - Struct containing details about the classes to be scheduled.
//
// @return map[string][]api.Class - Map where keys represent year and month, and values represent scheduled classes.
func separateClassByYearMonth(base api.ClassScheduler) map[string][]api.Class {
	datesMap := make(map[string][]api.Class)

	current := base.StartDate
	for current.Before(base.EndDate) || current.Equal(base.EndDate) {
		key := fmt.Sprintf("%d-%02d", current.Year(), current.Month())
		datesMap[key] = append(datesMap[key], api.Class{
			Name:     base.Name,
			Date:     current,
			Capacity: base.Capacity,
		})
		current = current.AddDate(0, 0, 1)
	}

	return datesMap
}
