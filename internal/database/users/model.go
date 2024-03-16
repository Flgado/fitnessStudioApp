package users

import "time"

type UserRow struct {
	Id             int       `db:"id"`
	Name           string    `db:"user_name"`
	CreateDate     time.Time `db:"create_date"`
	LastUpdateDate time.Time `db:"last_update_date"`
}
