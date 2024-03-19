package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
// @Param className query string false "Filter by class name"
// @Param startDate query string false "Filter classes with start date greater than or equal to the specified date. Format: dddd-dd-dd"
// @Param endDate query string false "Filter classes with end date less than or equal to the specified date. Format: dddd-dd-dd"
// @Param capacityGte query integer false "Filter classes with capacity greater than or equal to the specified value"
// @Param capacityLe query integer false "Filter classes with capacity less than or equal to the specified value"
// @Param numRegistrationsGte query integer false "Filter classes with number of registrations greater than or equal to the specified value"
// @Param numRegistrationsLe query integer false "Filter classes with number of registrations less than or equal to the specified value"
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
		return
	}

	result, err := h.uc.GetFilteredClasses(ctx, filters)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, result)
}

// HandlerAddClass handles the HTTP request to add a new class.
// @Summary Create multiple classes.
// @Description Creates new classes using the provided details. New classes will be created for each day within the range specified by the start date and end date.
// @Description If any of these days are unavailable, the endpoint will return the corresponding classes, indicating that scheduling was not possible
// @Tags Classes
// @Accept json
// @Produce json
// @Param body body api.ClassScheduler true "Class details (all fields are required, dates in the format YYYY-MM-DD)"
// @Success 200 {string} map[string]interface{}{"message": "All Classes Created With Success", "Not Possible To Scheduler": array<api.Class>}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/fitnessstudio/classes [post]
func (h ClassesHandler) HandlerAddClass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var addClass api.ClassSchedulerReceiver
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

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, addClass.StartDate)

	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	endDate, err := time.Parse(layout, addClass.EndDate)

	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	createClass := api.ClassScheduler{
		Name:      addClass.Name,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  addClass.Capacity,
	}
	// returned classes that was not possible to sheduler
	c, err := h.uc.CreateClass(ctx, createClass)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	if c != nil {
		respondWithJson(w, http.StatusOK, map[string][]api.Class{"Not Possible To Schedule": c})
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"message": "All Classes Created With Success"})
}

// HandlerUpdateClass handles the HTTP request to update a class.
// @Description Update class.
// @Tags Classes
// @Produce json
// @Param request body api.PatchClass{} true "Class data to update"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/classes [patch]
func (h ClassesHandler) HandlerUpdateClass(w http.ResponseWriter, r *http.Request) {

	patchClass := api.PatchClass{}

	err := json.NewDecoder(r.Body).Decode(&patchClass)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	updateClass, err := BuildUpdateClass(patchClass.Date, patchClass.Name, patchClass.Capacity)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}
	_, err = h.uc.UpdateClass(r.Context(), updateClass, patchClass.Id)
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
// @Param classId path int true "Class ID"
// @Success 200 {object} api.ReadClass
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/classes/{classId} [get]
func (h ClassesHandler) HandlerGetClassById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	classIdStr := chi.URLParam(r, "classId")

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
