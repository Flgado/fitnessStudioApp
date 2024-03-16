package usecases

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
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

func TestClassesUseCases_CreateClass(t *testing.T) {
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

	// Print the array
	for _, item := range mockWriteRepo.values {
		t.Log(item.Name)
		t.Log(item.Date)
	}
}
