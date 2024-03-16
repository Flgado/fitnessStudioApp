package api

import "time"

type ClassScheduler struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
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
	Name     *string    `json:"name"`
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
