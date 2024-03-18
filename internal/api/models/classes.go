package api

import "time"

type ClassScheduler struct {
	Name      string    `json:"name" validate:"required,len=1,max=50"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtefield=StartDate"`
	Capacity  int       `json:"capacity"`
}

type ReadClass struct {
	Id int `json:"id,omitempty"`
	Class
	NumRegistrations int `json:"num_registrations,omitempty"`
}

type Class struct {
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Capacity int       `json:"capacity"`
}

type UpdateClass struct {
	Name     *string    `json:"name" validate:"len=1,max=50"`
	Date     *time.Time `json:"date"`
	Capacity *int       `json:"capacity"`
}

type ClasseFilters struct {
	Name                string
	StartDateGte        *time.Time
	EndDateLe           *time.Time
	CapacityGte         *int
	CapacityLe          *int
	NumRegistrationsGte *int
	NumRegistrationsLe  *int
}
