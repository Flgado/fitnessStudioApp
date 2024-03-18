package handlers

import (
	"encoding/json"
	"net/http"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
)

type MakeReservationHandler struct {
	uc usecases.MakeBookUseCase
}

func NewMakeReservationHandler(uc usecases.MakeBookUseCase) *MakeReservationHandler {
	return &MakeReservationHandler{uc: uc}
}

// HandlerCreateBooking handles the HTTP request make a class reservation.
// @Description Make class reservation
// @Tags Bookings
// @Produce json
// @Param request body api.MakeRegervation true "Reservation body"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/bookings [post]
func (h *MakeReservationHandler) HandlerCreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reservation api.MakeRegervation
	err := json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	err = h.uc.Book(ctx, reservation.UserId, reservation.ClassId)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"message": "Succesfull Regervation"})
}
