package api

// @Description UserModel
type User struct {
	Id   int
	Name string
} //@name User model

// @Description UpdateUser Information
type UpdateUser struct {
	Name string `json:"name"`
} //@name Update User

// @Description CreateUser
type CreateUser struct {
	Name string `json:"name"`
} //@name CreateUser
