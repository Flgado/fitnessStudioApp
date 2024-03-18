package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/go-chi/chi"
)

type ClassesHandler struct {
	uc usecases.ClassesUseCases
}

func NewClassesHandler(uc usecases.ClassesUseCases) *ClassesHandler {
	return &ClassesHandler{uc: uc}
}

// HandlerGetClasses handles the HTTP request to get classes with optional filters.
// @Summary Get classes with optional filters
// @Description Returns a list of classes, optionally filtered by various parameters. If no filters are passed, it returns all classes.
// @Tags Classes
// @Produce json
// @Param class_name query string false "Filter by class name"
// @Param start_date query string false "Filter classes with start date greater than or equal to the specified date. Format: RFC3339"
// @Param end_date query string false "Filter classes with end date less than or equal to the specified date. Format: RFC3339"
// @Param capacity_gte query integer false "Filter classes with capacity greater than or equal to the specified value"
// @Param capacity_le query integer false "Filter classes with capacity less than or equal to the specified value"
// @Param num_registrations_gte query integer false "Filter classes with number of registrations greater than or equal to the specified value"
// @Param num_registrations_le query integer false "Filter classes with number of registrations less than or equal to the specified value"
// @Success 200 {array} api.ReadClass "Successful operation"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/classes [get]
func (h ClassesHandler) HandlerGetClasses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()

	// Create filter object
	filters, err := buildClassFilters(queryParams)
	if err != nil {
		responseWithErrors(w, *r, err)
	}

	result, err := h.uc.GetFilteredClasses(ctx, filters)

	if err != nil {
		responseWithErrors(w, *r, err)
	}

	respondWithJson(w, 200, result)
}

// HandlerAddClass handles the HTTP request to add a new class.
// @Summary Create multiple classes.
// @Description Adds new classes with the provided details.
// @Tags Classes
// @Accept json
// @Produce json
// @Param body body api.ClassScheduler true "Class details (all fields are required)"
// @Success 200 {string} map[string]interface{}{"message": "All Classes Created With Success", "Not Possible To Scheduler": array<api.Class>}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/fitnessstudio/classes [post]
func (h ClassesHandler) HandlerAddClass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var addClass api.ClassScheduler
	err := json.NewDecoder(r.Body).Decode(&addClass)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	// returned classes that was not possible to sheduler
	c, err := h.uc.CreateClass(ctx, addClass)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	if c != nil {
		respondWithJson(w, http.StatusOK, map[string][]api.Class{"Not Possible To Scheduler": c})
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"message": "All Classes Created With Success"})
}

// HandlerUpdateClass handles the HTTP request to update a class.
// @Description Update class.
// @Tags Classes
// @Produce json
// @Param class-id path int true "Class ID"
// @Param request body api.UpdateClass true "Class data to update"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/classes/{class-id} [post]
func (h ClassesHandler) HandlerUpdateClass(w http.ResponseWriter, r *http.Request) {
	classIdStr := chi.URLParam(r, "class-id")
	classId, err := strconv.Atoi(classIdStr)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Id with invalid format",
			"class id should be an integer")

		responseWithErrors(w, *r, e)
		return
	}

	var updateClass api.UpdateClass
	err = json.NewDecoder(r.Body).Decode(&updateClass)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	_, err = h.uc.UpdateClass(r.Context(), updateClass, classId)
	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"message": "Class Succesfull updated"})
}

// HandlerGetClassById handles the HTTP request to get a class by ID.
// @Description Get a class by ID
// @Tags Classes
// @Produce json
// @Param class-id path int true "Class ID"
// @Success 200 {object} api.ReadClass
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/classes/{class-id} [get]
func (h ClassesHandler) HandlerGetClassById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	classIdStr := chi.URLParam(r, "class-id")

	classId, err := strconv.Atoi(classIdStr)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	class, err := h.uc.GetClassById(ctx, classId)
	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, class)
}
