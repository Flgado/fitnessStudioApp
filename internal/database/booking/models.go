package booking

import "time"

type BokingRow struct {
	UserId       int       `db:"user_id"`
	ClassId      int       `db:"class_id"`
	ReservedDate time.Time `db:"reserved_date"`
}

type ClassBookedRow struct {
	Id               int       `db:"id"`
	Name             string    `db:"class_name"`
	Date             time.Time `db:"class_date"`
	Capacity         int       `db:"class_capacity"`
	NumRegistrations int       `db:"num_registrations"`
	CreateDate       time.Time `db:"create_date"`
	LastUpdateDate   time.Time `db:"last_update_date"`
	ReservedDate     time.Time `db:"reserved_date"`
}

type UserBookedRow struct {
	Id   int    `db:"id"`
	Name string `db:"user_name"`
}
