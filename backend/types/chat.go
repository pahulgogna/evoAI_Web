package types

import "time"

const (
	RoleSystem    MessageRole = "system"
	RoleAssistant MessageRole = "assistant"
	RoleUser      MessageRole = "user"
)

type ChatStore interface {
	StoreMessage(userId string, message *StoreMessage) (int32, error)
	GetAllChats(userId string) (*[]Chat, error)
	NewChatWithMessage(userId string, message *StoreMessage) (*NewChatWithMessageResponse, error)
	GetAllChatMessages(chatId string, userId string) (*[]Message, error)
	DeleteMessage(userId string, messageId string) error
	DeleteChat(userId string, chatId string) error
}

type MessageRole string

var ValidMessageSenders = map[MessageRole]bool{RoleUser: true, RoleSystem: true}

type Message struct {
	ID        int32       `json:"id" db:"id"`
	ChatId    int         `json:"chat_id" db:"chat_id"`
	Role      MessageRole `json:"role" db:"role"`
	Content      string      `json:"content" db:"content"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

type Chat struct {
	ID        int16     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserId    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Messages  []Message `json:"messages"`
}

type StoreMessage struct {
	ChatId        int16       `json:"chat_id" db:"chat_id" validate:"required"`
	Role          MessageRole `json:"role" db:"role" validate:"required"`
	Content          string      `json:"content" db:"content" validate:"required"`
	CreateNewChat bool        `json:"create_new_chat"`
}

type NewChatWithMessageResponse struct {
	ChatId    int16 `json:"chat_id"`
	MessageId int32 `json:"message_id"`
}
