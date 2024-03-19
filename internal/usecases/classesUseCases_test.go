//go:build unittests
// +build unittests

package usecases

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWriteRepository is a mock implementation of WriteRepository for testing purposes.
type MockWriteRepository struct {
	values []api.Class
	sync.Mutex
}

// Add mocks the Add method of WriteRepository.
func (m *MockWriteRepository) Add(ctx context.Context, classes []api.Class) error {
	m.Lock()
	defer m.Unlock()
	m.values = append(m.values, classes...)
	return nil
}

func (m *MockWriteRepository) Update(ctx context.Context, classId int, classUpdate api.UpdateClass) (int64, error) {
	return 2, nil
}

type mockClassesReadRepository struct {
	mock.Mock
}
type ReadRepository interface {
	List(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error)
	GetById(ctx context.Context, classId int) (api.ReadClass, error)
	GetClassReservations(ctx context.Context, classId int) (int, error)
}

func (m *mockClassesReadRepository) GetClassReservations(ctx context.Context, classId int) (int, error) {
	return 0, nil
}

func (m *mockClassesReadRepository) List(ctx context.Context, filters api.ClasseFilters) ([]api.ReadClass, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]api.ReadClass), args.Error(1)
}

func (m *mockClassesReadRepository) GetById(ctx context.Context, classId int) (api.ReadClass, error) {
	args := m.Called(ctx, classId)
	return args.Get(0).(api.ReadClass), args.Error(1)
}

type mockClassesWriteRepository struct {
	mock.Mock
}

func (m *mockClassesWriteRepository) Add(ctx context.Context, classes []api.Class) error {
	args := m.Called(ctx, classes)
	return args.Error(0)
}

func (m *mockClassesWriteRepository) Update(ctx context.Context, classId int, updateClass api.UpdateClass) (int64, error) {
	args := m.Called(ctx, classId, updateClass)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateClass_Success(t *testing.T) {
	mockReadRepo := new(mockClassesReadRepository)
	mockWriteRepo := new(mockClassesWriteRepository)

	uc := NewClassesUseCases(mockReadRepo, mockWriteRepo)

	classScheduler := api.ClassScheduler{
		Name:      "Test Class",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 1),
		Capacity:  10,
	}

	mockWriteRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	classes, err := uc.CreateClass(context.Background(), classScheduler)

	assert.NoError(t, err)
	assert.Empty(t, classes)
	mockWriteRepo.AssertExpectations(t)
}

func TestCreateClass_Error(t *testing.T) {
	mockReadRepo := new(mockClassesReadRepository)
	mockWriteRepo := new(mockClassesWriteRepository)

	uc := NewClassesUseCases(mockReadRepo, mockWriteRepo)

	classScheduler := api.ClassScheduler{
		Name:      "Test Class",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 1),
		Capacity:  10,
	}

	mockWriteRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("error adding class"))

	classes, err := uc.CreateClass(context.Background(), classScheduler)

	expectedClasses := []api.Class{
		{
			Name:     classScheduler.Name,
			Date:     classScheduler.StartDate,
			Capacity: classScheduler.Capacity,
		},
		{
			Name:     classScheduler.Name,
			Date:     classScheduler.EndDate,
			Capacity: classScheduler.Capacity,
		},
	}
	// Sort the expected and actual classes by date
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].Date.Before(classes[j].Date)
	})
	sort.Slice(expectedClasses, func(i, j int) bool {
		return expectedClasses[i].Date.Before(expectedClasses[j].Date)
	})

	assert.EqualError(t, err, "error adding class")
	mockWriteRepo.AssertExpectations(t)
	assert.Equal(t, expectedClasses[0].Name, classes[0].Name)
	compareDates(t, expectedClasses[0].Date, classes[0].Date)

	// Validate that the sorted classes match the expected values for Name and Date fields
	assert.Equal(t, expectedClasses[1].Name, classes[1].Name)
	assert.Equal(t, expectedClasses[1].Name, classes[1].Name)
}

func compareDates(t *testing.T, expected, actual time.Time) {
	expectedDate := time.Date(expected.Year(), expected.Month(), expected.Day(), 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(actual.Year(), actual.Month(), actual.Day(), 0, 0, 0, 0, time.UTC)
	if !expectedDate.Equal(actualDate) {
		t.Errorf("expected date %v, got %v", expectedDate, actualDate)
	}
}

func TestCreateClass_FoundInCacheUnavailabeDay(t *testing.T) {
	mockReadRepo := new(mockClassesReadRepository)
	mockWriteRepo := new(mockClassesWriteRepository)

	uc := NewClassesUseCases(mockReadRepo, mockWriteRepo)

	classScheduler := api.ClassScheduler{
		Name:      "Test Class",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 1),
		Capacity:  10,
	}

	mockWriteRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	// will find all days available in cache
	emptyList, err := uc.CreateClass(context.Background(), classScheduler)

	classScheduler2 := api.ClassScheduler{
		Name:      "Not found available day",
		StartDate: time.Now(),
		EndDate:   time.Now(),
		Capacity:  10,
	}

	class, err2 := uc.CreateClass(context.Background(), classScheduler2)

	assert.Empty(t, emptyList)
	assert.Nil(t, err)
	assert.Nil(t, err2)
	assert.Len(t, class, 1)
	assert.Equal(t, classScheduler2.Name, class[0].Name)
	assert.Equal(t, classScheduler2.StartDate.Day(), class[0].Date.Day())
}

func TestClassesUseCases_CreateClassThreadSafe(t *testing.T) {
	// Initialize your use case with a mock WriteRepository
	mockWriteRepo := &MockWriteRepository{}
	useCase := classesUseCases{
		wrRep:        mockWriteRepo,
		reservedDays: sync.Map{},
	}

	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 10), Capacity: 10},
		{Name: "Test2", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 5), Capacity: 10},
		{Name: "Test3", StartDate: time.Now().AddDate(0, 0, 5), EndDate: time.Now().AddDate(0, 0, 15), Capacity: 10},
		{Name: "Test4", StartDate: time.Now().AddDate(0, 0, 15), EndDate: time.Now().AddDate(0, 0, 20), Capacity: 10},
		{Name: "Test5", StartDate: time.Now().AddDate(0, 1, 0), EndDate: time.Now().AddDate(0, 1, 5), Capacity: 10},
		{Name: "Test6", StartDate: time.Now().AddDate(0, 1, 0), EndDate: time.Now().AddDate(0, 1, 5), Capacity: 10},
		{Name: "Test7", StartDate: time.Now().AddDate(0, 1, 0), EndDate: time.Now().AddDate(0, 1, 10), Capacity: 10},
		{Name: "Test8", StartDate: time.Now(), EndDate: time.Now(), Capacity: 10},
		{Name: "Test9", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 3), Capacity: 10},
		{Name: "Test9", StartDate: time.Now().AddDate(0, 1, 0), EndDate: time.Now().AddDate(0, 1, 0), Capacity: 10},
	}

	// Define how many goroutines to run concurrently
	numGoroutines := 10

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run multiple goroutines concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			// Call CreateClass with each classScheduler in the table test
			_, err := useCase.CreateClass(context.Background(), data[index])
			if err != nil {
				t.Errorf("Error creating class: %v", err)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Use a map to store unique dates encountered
	uniqueDates := make(map[string]struct{})

	// Iterate through mockWriteRepo.values and validate dates
	for _, class := range mockWriteRepo.values {
		y, m, d := class.Date.Date()
		// Extract the date part without hours, minutes, and seconds
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)

		// Check if the date already exists in the map
		if _, ok := uniqueDates[date]; ok {
			t.Errorf("Duplicate date found: %v", date)
			break // Exit loop on first duplicate found
		}

		// Add the date to the map
		uniqueDates[date] = struct{}{}
	}
}
