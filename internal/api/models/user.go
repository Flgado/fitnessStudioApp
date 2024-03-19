package api

type UsersList struct {
	Users []User `json:"users,omitempty"`
}

// @Description UserModel
type User struct {
	Id   int    `json:"id,omitempty" validate:"required,len=1,max=50"`
	Name string `json:"name,omitempty"`
} //@name User model

// @Description UpdateUser Information
type UpdateUser struct {
	Name string `json:"name" validate:"required,len=1,max=50"`
} //@name Update User

// @Description CreateUser
type CreateUser struct {
	Name string `json:"name" validate:"required,len=1,max=50"`
} //@name CreateUser
