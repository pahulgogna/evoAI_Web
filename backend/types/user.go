package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	RemoveUser(id string) error
	CreateUser(user *User) error
}


type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}


type RegisterUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResponseUser struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Email   string `json:"email" db:"email"`
	IsAdmin bool   `json:"is_admin" db:"is_admin"`
}

type JWTUser struct {
	Id      string
	IsAdmin bool
}
