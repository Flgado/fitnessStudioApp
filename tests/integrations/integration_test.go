//go:build integrationtests
// +build integrationtests

package integrations

import (
	"context"
	"log"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/internal/database/booking"
	"github.com/Flgado/fitnessStudioApp/internal/database/classes"
	"github.com/Flgado/fitnessStudioApp/internal/database/users"
	"github.com/Flgado/fitnessStudioApp/internal/usecases"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var testDbInstance *sqlx.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func cleanupUserTableDatabase() {
	_, err := testDbInstance.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Error cleaning up database: %v", err)
	}
}

func resetAutoIncrement() {
	_, err := testDbInstance.Exec("ALTER SEQUENCE classes_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatalf("Error resetting auto-increment counter: %v", err)
	}
}

func cleanupClassesTableDatabase() {
	_, err := testDbInstance.Exec("DELETE FROM classes")
	if err != nil {
		log.Fatalf("Error cleaning up database: %v", err)
	}

	resetAutoIncrement()
}

func cleanupAllTablesDatabase() {
	_, err := testDbInstance.Exec("DELETE FROM booking")
	if err != nil {
		log.Fatalf("Error cleaning up database: %v", err)
	}
	cleanupClassesTableDatabase()
	cleanupUserTableDatabase()
	_, err = testDbInstance.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatalf("Error resetting auto-increment counter: %v", err)
	}

}

func TestCreateUser(t *testing.T) {
	defer cleanupUserTableDatabase()
	// Arrange
	wriRep := users.NewReadRepository(testDbInstance)
	readRep := users.NewWriteRepository(testDbInstance)

	expectedUser := api.User{
		Id:   1,
		Name: "Joao Folgado",
	}

	ctx := context.Background()
	uc := usecases.NewUserUseCase(wriRep, readRep)

	// act
	err1 := uc.CreateUser(ctx, "Joao Folgado")
	user, err2 := uc.GetUserById(ctx, 1)

	// assert
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, user, expectedUser)
}

func TestGetAllUsers(t *testing.T) {
	defer cleanupUserTableDatabase()
	// Arrange
	wriRep := users.NewReadRepository(testDbInstance)
	readRep := users.NewWriteRepository(testDbInstance)

	expectedUser := []api.User{
		{Id: 1,
			Name: "Joao Folgado"},
		{
			Id:   2,
			Name: "Joao Folgado2",
		},
	}

	ctx := context.Background()
	uc := usecases.NewUserUseCase(wriRep, readRep)

	// act
	err1 := uc.CreateUser(ctx, "Joao Folgado1")
	err2 := uc.CreateUser(ctx, "Joao Folgado2")
	users, err3 := uc.GetAllUsers(ctx)

	// assert
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUser, expectedUser)
}

func TestCreateClasses(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	startDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	expectedClasses := []api.ReadClass{
		{Id: 1, Class: api.Class{Name: "Test1", Date: startDate, Capacity: 10}, NumRegistrations: 0},
		{Id: 2, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 1), Capacity: 10}, NumRegistrations: 0},
		{Id: 3, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 2), Capacity: 10}, NumRegistrations: 0},
		{Id: 4, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 3), Capacity: 10}, NumRegistrations: 0},
		{Id: 5, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 4), Capacity: 10}, NumRegistrations: 0},
		{Id: 6, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 5), Capacity: 10}, NumRegistrations: 0},
	}

	sort.Slice(expectedClasses, func(i, j int) bool {
		return expectedClasses[i].Id < expectedClasses[j].Id
	})

	filters := api.ClasseFilters{
		Name: "Test1",
	}

	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	noPossibleToScheduler, err1 := uc.CreateClass(ctx, data[0])
	allClasses, err2 := uc.GetFilteredClasses(ctx, filters)

	// assert
	assert.Nil(t, noPossibleToScheduler)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Len(t, expectedClasses, 6)
	assert.Equal(t, expectedClasses, allClasses)
}

func TestCreateClassesForUnavailableDays(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	startDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	willFoundNotAvailableDays := api.ClassScheduler{
		Name:      "Test2",
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 5),
		Capacity:  10,
	}

	expectedClasses := []api.ReadClass{
		{Id: 1, Class: api.Class{Name: "Test1", Date: startDate, Capacity: 10}, NumRegistrations: 0},
		{Id: 2, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 1), Capacity: 10}, NumRegistrations: 0},
		{Id: 3, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 2), Capacity: 10}, NumRegistrations: 0},
		{Id: 4, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 3), Capacity: 10}, NumRegistrations: 0},
		{Id: 5, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 4), Capacity: 10}, NumRegistrations: 0},
		{Id: 6, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 5), Capacity: 10}, NumRegistrations: 0},
	}

	notPossibleToAddExpected := []api.Class{
		{Name: "Test2", Date: startDate, Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 1), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 2), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 3), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 4), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	sort.Slice(expectedClasses, func(i, j int) bool {
		return expectedClasses[i].Id < expectedClasses[j].Id
	})

	filters := api.ClasseFilters{
		Name: "Test1",
	}

	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	noPossibleToScheduler, err1 := uc.CreateClass(ctx, data[0])
	noPossibleToScheduler2, err3 := uc.CreateClass(ctx, willFoundNotAvailableDays)
	allClasses, err2 := uc.GetFilteredClasses(ctx, filters)

	// assert
	assert.Nil(t, noPossibleToScheduler)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Len(t, expectedClasses, 6)
	assert.Len(t, noPossibleToScheduler2, 6)
	assert.Equal(t, noPossibleToScheduler2, notPossibleToAddExpected)
	assert.Equal(t, expectedClasses, allClasses)
}

func TestCreateClassesForParcialNotAvailableDays(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	startDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	willFoundNotAvailableDays := api.ClassScheduler{
		Name:      "Test2",
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 0, 8),
		Capacity:  10,
	}

	expectedClasses := []api.ReadClass{
		{Id: 1, Class: api.Class{Name: "Test1", Date: startDate, Capacity: 10}, NumRegistrations: 0},
		{Id: 2, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 1), Capacity: 10}, NumRegistrations: 0},
		{Id: 3, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 2), Capacity: 10}, NumRegistrations: 0},
		{Id: 4, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 3), Capacity: 10}, NumRegistrations: 0},
		{Id: 5, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 4), Capacity: 10}, NumRegistrations: 0},
		{Id: 6, Class: api.Class{Name: "Test1", Date: startDate.AddDate(0, 0, 5), Capacity: 10}, NumRegistrations: 0},
		{Id: 7, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 6), Capacity: 10}, NumRegistrations: 0},
		{Id: 8, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 7), Capacity: 10}, NumRegistrations: 0},
		{Id: 9, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 8), Capacity: 10}, NumRegistrations: 0},
	}

	notPossibleToAddExpected := []api.Class{
		{Name: "Test2", Date: startDate, Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 1), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 2), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 3), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 4), Capacity: 10},
		{Name: "Test2", Date: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	sort.Slice(expectedClasses, func(i, j int) bool {
		return expectedClasses[i].Id < expectedClasses[j].Id
	})

	filters := api.ClasseFilters{
		Name: "",
	}

	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	noPossibleToScheduler, err1 := uc.CreateClass(ctx, data[0])
	noPossibleToScheduler2, err3 := uc.CreateClass(ctx, willFoundNotAvailableDays)
	allClasses, err2 := uc.GetFilteredClasses(ctx, filters)

	// assert
	assert.Nil(t, noPossibleToScheduler)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Len(t, expectedClasses, 9)
	assert.Len(t, noPossibleToScheduler2, 6)
	assert.Equal(t, noPossibleToScheduler2, notPossibleToAddExpected)
	assert.Equal(t, expectedClasses, allClasses)
}

func TestGetClassesWithFilters(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	startDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 5},
		{Name: "Test2", StartDate: startDate.AddDate(0, 0, 6), EndDate: startDate.AddDate(0, 0, 10), Capacity: 10},
	}

	filter1 := api.ClasseFilters{
		Name: "Test2",
	}

	expectedForFilter1 := []api.ReadClass{
		{Id: 7, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 6), Capacity: 10}, NumRegistrations: 0},
		{Id: 8, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 7), Capacity: 10}, NumRegistrations: 0},
		{Id: 9, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 8), Capacity: 10}, NumRegistrations: 0},
		{Id: 10, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 9), Capacity: 10}, NumRegistrations: 0},
		{Id: 11, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 10), Capacity: 10}, NumRegistrations: 0},
	}

	dateFilter2 := startDate.AddDate(0, 0, 8)
	filter2 := api.ClasseFilters{
		Name:         "Test2",
		StartDateGte: &dateFilter2,
	}

	expectedForFilter2 := []api.ReadClass{
		{Id: 9, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 8), Capacity: 10}, NumRegistrations: 0},
		{Id: 10, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 9), Capacity: 10}, NumRegistrations: 0},
		{Id: 11, Class: api.Class{Name: "Test2", Date: startDate.AddDate(0, 0, 10), Capacity: 10}, NumRegistrations: 0},
	}

	filter3 := api.ClasseFilters{
		Name:        "",
		CapacityGte: Int(6),
	}

	testCases := []struct {
		testName string
		filter   api.ClasseFilters
		resutls  []api.ReadClass
	}{
		{"filter1", filter1, expectedForFilter1},
		{"filter2", filter2, expectedForFilter2},
		{"filter3", filter3, expectedForFilter1},
	}
	// act
	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)
	for _, d := range data {
		_, _ = uc.CreateClass(ctx, d)
	}

	for _, tc := range testCases {
		classes, err := uc.GetFilteredClasses(ctx, tc.filter)

		// Assert that there's no error
		assert.Nil(t, err)

		// Assert that the classes obtained are equal to the expected results
		assert.Equal(t, tc.resutls, classes, "Unexpected classes result for test case: %s", tc.testName)
	}
}

func TestUpdateExistingClassDate_DateUnavailable(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	y, m, d := time.Now().Date()
	startDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	updateDate := startDate.AddDate(0, 0, 3)
	updateClass := api.UpdateClass{
		Name: String("Updated"),
		Date: &updateDate,
	}

	class := api.Class{
		Name:     "Updated",
		Date:     updateDate,
		Capacity: 20,
	}

	expectedError := utils.E(http.StatusNotFound,
		nil,
		map[string]string{"message": "Date already reserved"},
		"The selected date is already reserved.",
		"Please choose a different date or class.")

	expectedClass := api.ReadClass{
		Id:               5,
		Class:            class,
		NumRegistrations: 0,
	}
	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	_, err := uc.CreateClass(ctx, data[0])
	rowsUpdated, err1 := uc.UpdateClass(ctx, updateClass, 5)
	classToValidate, err2 := uc.GetClassById(ctx, 5)
	// assert

	assert.Nil(t, err)
	assert.Equal(t, expectedError, err1)
	assert.Nil(t, err2)
	assert.Equal(t, rowsUpdated, int64(0))
	assert.NotEqual(t, expectedClass, classToValidate)
}

func TestUpdateExistingClassDate(t *testing.T) {
	defer cleanupClassesTableDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)

	y, m, d := time.Now().Date()
	startDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	data := []api.ClassScheduler{
		{Name: "Test1", StartDate: startDate, EndDate: startDate.AddDate(0, 0, 5), Capacity: 10},
	}

	updateDate := startDate.AddDate(0, 0, 8)
	updateClass := api.UpdateClass{
		Name: String("Updated"),
		Date: &updateDate,
	}

	class := api.Class{
		Name:     "Updated",
		Date:     updateDate,
		Capacity: 10,
	}

	expectedClass := api.ReadClass{
		Id:               5,
		Class:            class,
		NumRegistrations: 0,
	}
	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	_, err := uc.CreateClass(ctx, data[0])
	rowsUpdated, err1 := uc.UpdateClass(ctx, updateClass, 5)
	classToValidate, err2 := uc.GetClassById(ctx, 5)
	// assert

	assert.Nil(t, err)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, rowsUpdated, int64(1))
	assert.Equal(t, expectedClass, classToValidate)
}

func TestUpdateExisting_NewCapacityLowerThenNumRegistrations(t *testing.T) {
	defer cleanupAllTablesDatabase()
	// Arrange
	wriRep := classes.NewReadRepository(testDbInstance)
	readRep := classes.NewWriteRepository(testDbInstance)
	date := "2024-03-17"
	testDbInstance.DB.Exec(`INSERT INTO classes (class_name, class_date, class_capacity, num_registrations)
	VALUES('Test', '` + date + `', 3, 3)`)

	updateClass := api.UpdateClass{
		Name:     String("Updated"),
		Capacity: Int(2),
	}

	layout := "2006-01-02"
	expectedDate, err := time.Parse(layout, date)

	if err != nil {
		t.Fatal("Not possible to parse date")
	}

	class := api.Class{
		Name:     "Test",
		Date:     expectedDate,
		Capacity: 3,
	}

	expectedError := utils.E(http.StatusUnprocessableEntity,
		nil,
		map[string]string{"message": "Cannot update class capacity"},
		"The class is on full capacity.",
		"Please try again later or select a different class.")

	expectedClass := api.ReadClass{
		Id:               1,
		Class:            class,
		NumRegistrations: 3,
	}
	ctx := context.Background()
	uc := usecases.NewClassesUseCases(wriRep, readRep)

	// act
	rowsUpdated, err := uc.UpdateClass(ctx, updateClass, 1)
	classToValidate, err1 := uc.GetClassById(ctx, 1)
	// assert

	assert.Equal(t, expectedError, err)
	assert.Nil(t, err1)
	assert.Equal(t, rowsUpdated, int64(0))
	assert.Equal(t, expectedClass, classToValidate)
}

func TestCreateReservation_ClassAlreadyFull(t *testing.T) {
	defer cleanupAllTablesDatabase()
	// Arrange
	date := "2024-03-17"
	testDbInstance.DB.Exec(`INSERT INTO classes (class_name, class_date, class_capacity, num_registrations)
	VALUES('Test', '` + date + `', 3, 3)`)
	testDbInstance.DB.Exec(`INSERT INTO users (user_name) VALUES('Joao Folgado')`)

	redRep := booking.NewReadRepository(testDbInstance)
	wrRep := booking.NewWriteRepository(testDbInstance)

	bookUseCase := usecases.NewBookUseCase(redRep, wrRep)
	makeReservationUseCase := usecases.NewMakeBookUseCase(redRep, wrRep)

	expectedError := utils.E(http.StatusUnprocessableEntity,
		nil,
		map[string]string{"message": "Class Capacity Reached"},
		"The class is already full and cannot accept any more registrations.",
		"Please select another class or try again later.")

	// Act
	err := makeReservationUseCase.Book(context.Background(), 1, 1)

	// assert
	assert.Equal(t, expectedError, err)
	classReservations, err2 := bookUseCase.GetClassesReservations(context.Background(), 1)
	assert.Len(t, classReservations, 0)
	assert.Nil(t, err2)
}

func TestCreateReservation_ClassAvailable(t *testing.T) {
	cleanupAllTablesDatabase()
	// Arrange
	date := "2024-03-17T12:00:00Z"
	testDbInstance.DB.Exec(`INSERT INTO classes (class_name, class_date, class_capacity, num_registrations)
	VALUES('Test', '` + date + `', 3, 2)`)
	testDbInstance.DB.Exec(`INSERT INTO users (user_name) VALUES('Joao Folgado')`)

	redRep := booking.NewReadRepository(testDbInstance)
	wrRep := booking.NewWriteRepository(testDbInstance)

	bookUseCase := usecases.NewBookUseCase(redRep, wrRep)
	makeReservationUseCase := usecases.NewMakeBookUseCase(redRep, wrRep)

	expectedResult := []api.UsersBooked{
		{ClassId: 1, UserId: 1, UserName: "Joao Folgado"},
	}
	// Act
	err := makeReservationUseCase.Book(context.Background(), 1, 1)

	// assert
	assert.Nil(t, err)
	classReservationsList, err2 := bookUseCase.GetClassesReservations(context.Background(), 1)
	usersReservationList, err3 := bookUseCase.GetUserReservations(context.Background(), 1)
	assert.Len(t, classReservationsList, 1)
	assert.Equal(t, expectedResult, classReservationsList)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Len(t, usersReservationList, 1)
}

func Int(i int) *int {
	return &i
}

func String(i string) *string {
	return &i
}
