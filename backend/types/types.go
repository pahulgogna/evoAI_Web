package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	RemoveUser(id string) error
	CreateUser(user *User) error
}

type ToolStore interface {
	GetToolByIdAndUserId(id string, userId string) (*Tool, error)
	GetToolsByUserId(userId string) (*[]Tool, error)
	CreateToolByUser(tool *CreateTool, userId string) error
	UpdateTool(tool *UpdateTool, id string, userId string) error
	DeleteTool(id string, userId string) error
}

type ChatStore interface {
	StoreMessage(userId string, messages *StoreMessage) (int32, error)
	GetAllChats(userId string) (*[]Chat, error)
	NewChatWithMessage(userId string, message *StoreMessage) (*NewChatWithMessageResponse, error)
	GetAllChatMessages(chatId string, userId string) (*[]Message, error)
	DeleteMessage(userId string, messageId string) error
	DeleteChat(userId string, chatId string) error
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

type UpdateTool struct {
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	Language     *string   `json:"language"`
	Code         *string   `json:"code"`
	Dependencies *[]string `json:"dependencies"`
}

type Message struct {
	ID        int32     `json:"id" db:"id"`
	ChatId    int       `json:"chat_id" db:"chat_id"`
	By        string    `json:"by" db:"by"`
	Data      string    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Chat struct {
	ID        int16     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserId    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Messages  []Message `json:"messages"`
}

type StoreMessage struct {
	ChatId        int16  `json:"chat_id" db:"chat_id" validate:"required"`
	By            string `json:"by" db:"by" validate:"required"`
	Data          string `json:"data" db:"data" validate:"required"`
	CreateNewChat bool   `json:"create_new_chat"`
}

type NewChatWithMessageResponse struct {
	ChatId    int16 `json:"chat_id"`
	MessageId int32 `json:"message_id"`
}
