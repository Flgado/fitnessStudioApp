//go:build unittests
// +build unittests

package usecases_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBookingReadRepository struct {
	mock.Mock
}

func (m *mockBookingReadRepository) GetUserBookings(ctx context.Context, userId int) ([]api.ClassBooked, error) {
	return nil, nil
}

func (m *mockBookingReadRepository) GetClassReservations(ctx context.Context, classId int) ([]api.UsersBooked, error) {
	return nil, nil
}

func (m *mockBookingReadRepository) IsClassBookedByUser(ctx context.Context, userId int, classId int) (bool, error) {
	args := m.Called(ctx, userId, classId)
	return args.Bool(0), args.Error(1)
}

type mockBookingWriteRepository struct {
	mock.Mock
}

func (m *mockBookingWriteRepository) Add(ctx context.Context, userId int, classId int) error {
	args := m.Called(ctx, userId, classId)
	return args.Error(0)
}

// TestBook_Success tests the Book method when the class is successfully booked
func TestBook_Success(t *testing.T) {
	// Initialize mock repositories
	mockReadRepo := new(mockBookingReadRepository)
	mockWriteRepo := new(mockBookingWriteRepository)

	// Create the use case with mock repositories
	uc := usecases.NewMakeBookUseCase(mockReadRepo, mockWriteRepo)

	// Define test data
	userId := 123
	classId := 456

	// Set up mock behavior
	mockReadRepo.On("IsClassBookedByUser", mock.Anything, userId, classId).Return(false, nil)
	mockWriteRepo.On("Add", mock.Anything, userId, classId).Return(nil)

	// Call the method under test
	err := uc.Book(context.Background(), userId, classId)

	// Assertions
	assert.NoError(t, err)
	mockReadRepo.AssertExpectations(t)
	mockWriteRepo.AssertExpectations(t)
}

// TestBook_Conflict tests the Book method when the class is already booked by the user
func TestBook_Conflict(t *testing.T) {
	// Initialize mock repositories
	mockReadRepo := new(mockBookingReadRepository)
	mockWriteRepo := new(mockBookingWriteRepository)

	// Create the use case with mock repositories
	uc := usecases.NewMakeBookUseCase(mockReadRepo, mockWriteRepo)

	// Define test data
	userId := 123
	classId := 456

	// Set up mock behavior
	mockReadRepo.On("IsClassBookedByUser", mock.Anything, userId, classId).Return(true, nil)

	// Call the method under test
	err := uc.Book(context.Background(), userId, classId)

	// Assertions
	assert.EqualError(t, err, utils.E(http.StatusConflict, nil, map[string]string{"message": "Conflict Status"},
		"Class with Id: 123 is already reserved by User with id 456", "Validate user reserved classes").Error())

	mockReadRepo.AssertExpectations(t)
}

// TestBook_RepositoryError tests the Book method when there's an error with the repository
func TestBook_RepositoryError(t *testing.T) {
	// Initialize mock repositories
	mockReadRepo := new(mockBookingReadRepository)
	mockWriteRepo := new(mockBookingWriteRepository)

	// Create the use case with mock repositories
	uc := usecases.NewMakeBookUseCase(mockReadRepo, mockWriteRepo)

	// Define test data
	userId := 123
	classId := 456

	// Set up mock behavior
	mockReadRepo.On("IsClassBookedByUser", mock.Anything, userId, classId).Return(false, errors.New("repository error"))

	// Call the method under test
	err := uc.Book(context.Background(), userId, classId)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "repository error")

	mockReadRepo.AssertExpectations(t)
}
