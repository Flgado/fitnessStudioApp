package classes

import "time"

type ClassRow struct {
	Id               int       `db:"id"`
	Name             string    `db:"class_name"`
	Date             time.Time `db:"class_date"`
	Capacity         int       `db:"class_capacity"`
	NumRegistrations int       `db:"num_registrations"`
	CreateDate       time.Time `db:"create_date"`
	LastUpdateDate   time.Time `db:"last_update_date"`
}
