package api

import "time"

type ClassBooked struct {
	Id           int       `json:"class_id,omitempty"`
	Name         string    `json:"class_name,omitempty"`
	Date         time.Time `json:"class_date,omitempty"`
	ReservedDate time.Time `json:"reserved_date,omitempty"`
}

type UsersBooked struct {
	ClassId  int    `json:"class_id,omitempty"`
	UserId   int    `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

type MakeBooking struct {
	ClassId int `json:"class_id,omitempty"`
	UserId  int `json:"user_id,omitempty"`
}
