package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/go-chi/chi"
)

type UsersHandler struct {
	uc usecases.UserUseCases
}

func NewUsersHandler(uc usecases.UserUseCases) *UsersHandler {
	return &UsersHandler{uc: uc}
}

type GetAllUsers struct {
	Users []api.UpdateUser
} // @name GetAllUsers

// HandlerGetUsers handles the HTTP request to get all users.
// @Description Get all users
// @Tags Users
// @Produce json
// @Success 200 {object} []User
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/users [get]
func (h UsersHandler) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.uc.GetAllUsers(ctx)
	if err != nil {
		responseWithError(w, 500, err.Error())
	}
	respondWithJson(w, 200, users)
}

// HandlerGetUserById handles the HTTP request to get a user by ID.
// @Description Get a user by ID
// @Tags Users
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/users/{userId} [get]
func (h UsersHandler) HandlerGetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdStr := chi.URLParam(r, "userId")
	// Convert the userId string to an integer
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	user, err := h.uc.GetUserById(ctx, userId)
	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, user)
}

// HandlerCreateUser handles the HTTP request to create a new user.
// @Description Create a new user.
// @Tags Users
// @Produce json
// @Param request body api.CreateUser true "User data to create"
// @Success 200  {object} api.CreateUser
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/users [post]
func (h UsersHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user api.CreateUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	if user.Name == "" {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"User name should not be empty",
			"Use a valid user name")

		responseWithErrors(w, *r, e)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	err = h.uc.CreateUser(r.Context(), user.Name)

	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, 200, map[string]string{"message": "User created with Success"})
}

// HandlerUpdateUser handles the HTTP request to update a user.
// @Description Update a user
// @Tags Users
// @Produce json
// @Param request body api.User true "User data to update"
// @Success 200
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/fitnessstudio/users [patch]
func (h UsersHandler) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := api.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"Request body not expected",
			"Read our documentation for more details")

		responseWithErrors(w, *r, e)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		e := utils.E(http.StatusBadRequest,
			err,
			map[string]string{"message": "BadRequest"},
			"User name should not be empty",
			"Use a valid user name")

		responseWithErrors(w, *r, e)
		return
	}
	_, err = h.uc.UpdateUser(r.Context(), user)
	if err != nil {
		responseWithErrors(w, *r, err)
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"Success": "Updated user"})
}
