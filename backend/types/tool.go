package types

import "time"


type ToolStore interface {
	GetToolByIdAndUserId(id string, userId string) (*Tool, error)
	GetToolsByUserId(userId string) (*[]Tool, error)
	CreateToolByUser(tool *CreateTool, userId string) error
	UpdateTool(tool *UpdateTool, id string, userId string) error
	DeleteTool(id string, userId string) error
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

type DBTool struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Language     string    `json:"language" db:"language"`
	Code         string    `json:"code" db:"code"`
	Dependencies string    `json:"dependencies" db:"dependencies"`
	UserId       int       `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type CreateTool struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Language     string   `json:"language" validate:"required"`
	Code         string   `json:"code" validate:"required"`
	Dependencies []string `json:"dependencies"`
}

type UpdateTool struct {
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	Language     *string   `json:"language"`
	Code         *string   `json:"code"`
	Dependencies *[]string `json:"dependencies"`
}
