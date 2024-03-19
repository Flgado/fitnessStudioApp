package handlers

import (
	"net/http"
	"strconv"

	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/go-chi/chi"
)

type BookingInfoHandler struct {
	uc usecases.BookingUseCase
}

func NewBookingInfoHandler(uc usecases.BookingUseCase) *BookingInfoHandler {
	return &BookingInfoHandler{uc: uc}
}

// HandlerGetUserClasses handles the HTTP request to get classe booked by user
// @Summary Get the list of classes by user
// @Description Returns a list of classes booked by user
// @Tags Bookings
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} []ClassBooked
// @Failure 400 {object} ErrorResponse
// @Router /v1/fitnessstudio/bookings/users/{userId}/classes [get]
func (h BookingInfoHandler) HandlerGetUserClasses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdStr := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"UserId should be integer",
			"Read our documentation")

		responseWithErrors(w, *r, e)
		return
	}

	result, err := h.uc.GetUserReservations(ctx, userId)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, result)
}

// HandlerGetClassUsers handles the HTTP request to get users registered in a class
// @Summary Get a list of users who have booked the class.
// @Description Return the list of users who have booked the class.
// @Tags Bookings
// @Produce json
// @Param classId path int true "Class Id"
// @Success 200 {array} []api.UsersBooked
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/bookings/classes/{classId}/users [get]
func (h BookingInfoHandler) HandlerGetClassUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	classIdStr := chi.URLParam(r, "classId")
	classId, err := strconv.Atoi(classIdStr)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"ClassId should be integer",
			"Read our documentation")

		responseWithErrors(w, *r, e)
		return
	}

	result, err := h.uc.GetClassesReservations(ctx, classId)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, result)
}
