package booking

const (
	AddBokking = `INSERT INTO bokking (user_id, class_id) 
					VALUES($1, $2)`

	GetUserBookings = `SELECT c.id, c.class_name, c.class_date, c.class_capacity, c.num_registrations, b.reserved_date
						FROM classes c
						INNER JOIN booking b ON c.id = b.class_id
						WHERE b.user_id = $1;
						`

	GetUsersOfBooking = `SELECT u.id, u.user_name
						FROM users u
						INNER JOIN booking b ON u.id = b.user_id
						WHERE b.class_id = $1;
						`
)
