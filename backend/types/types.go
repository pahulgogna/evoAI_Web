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

type Tool struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Language     string    `json:"language" db:"language"`
	Code         string    `json:"code" db:"code"`
	Dependencies []string  `json:"dependencies" db:"dependencies"`
	UserId       int       `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
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

type CreateTool struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Language     string   `json:"language" validate:"required"`
	Code         string   `json:"code" validate:"required"`
	Dependencies []string `json:"dependencies"`
}

