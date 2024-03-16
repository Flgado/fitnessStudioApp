package users

const (
	findUsers = `SELECT *
				  FROM users 
				  `
	findUserById = `SELECT *
						From users
						Where id = $1`

	findUserByName = `SELECT *
						From users
						Where name = ?`

	AddUserRow = `INSERT INTO users (user_name) VALUES(:user_name)`

	UpdateUser = `UPDATE users SET user_name =:user_name WHERE id =:id`
)
