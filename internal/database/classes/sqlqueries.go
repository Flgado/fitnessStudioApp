package classes

const (
	findClassById = `SELECT *
						From classes
						Where id = ?`

	classReservationsById = `SELECT num_registrations
								From classes
								Where id = ?`
	AddClassRow = `INSERT INTO classes (class_name, class_date, class_capacity, num_registrations) 
					VALUES(:class_name, :class_date, :class_capacity, :num_registrations)`

	UpdateClass = `UPDATE classes SET`
)
