package api

import "time"

type ClassSchedulerReceiver struct {
	Name      string `json:"name" validate:"required,len=1,max=50"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required,gtefield=StartDate"`
	Capacity  int    `json:"capacity"`
} //@name ClassSchedulerReceiver

type ClassScheduler struct {
	Name      string    `json:"name" validate:"required,len=1,max=50"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtefield=StartDate"`
	Capacity  int       `json:"capacity"`
} //@name ClassScheduler

type ReadClass struct {
	Id int `json:"id,omitempty"`
	Class
	NumRegistrations int `json:"num_registrations,omitempty"`
} // @name ReadClass

type Class struct {
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Capacity int       `json:"capacity"`
} // @name Class

type PatchClass struct {
	Id       int     `json:"id,omitempty"`
	Name     *string `json:"name,omitempty" validate:"len=1,max=50"`
	Date     *string `json:"date,omitempty"`
	Capacity *int    `json:"capacity,omitempty"`
} // @name PatchClass

type UpdateClass struct {
	Name     *string
	Date     *time.Time
	Capacity *int
} // @name UpdateClass

type ClasseFilters struct {
	Name                string
	StartDateGte        *time.Time
	EndDateLe           *time.Time
	CapacityGte         *int
	CapacityLe          *int
	NumRegistrationsGte *int
	NumRegistrationsLe  *int
} // @name ClasseFilters
