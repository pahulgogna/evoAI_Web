package chat

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pahulgogna/evoAI_Web/backend/types"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) StoreMessage(userId string, message *types.StoreMessage) (int32, error) {

	chatExistsCheckRows, err := s.db.Query("SELECT * FROM chat WHERE user_id = $1 AND id = $2;", userId, message.ChatId)
	if err != nil {
		return -1, err
	}
	if !chatExistsCheckRows.Next() {
		return -1, fmt.Errorf("chat not found")
	}

	rows, err := s.db.Query("INSERT INTO message (chat_id, by, data) VALUES ($1, $2, $3) RETURNING id;", message.ChatId, message.By, message.Data)
	if err != nil {
		return -1, err
	}
	var messageId int32
	for rows.Next() {
		rows.Scan(&messageId)
	}
	return messageId, nil
}

func (s *Store) NewChatWithMessage(userId string, message *types.StoreMessage) (*types.NewChatWithMessageResponse, error) {
	
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred in NewChatWithMessage: ", err)
		}
	}()

	tx := s.db.MustBegin()
	
	chatQueryRows, err := tx.Query("INSERT INTO chat (user_id) VALUES ($1) RETURNING id;", userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var chatId int16
	for chatQueryRows.Next() {
		chatQueryRows.Scan(&chatId)
	}

	messageQueryRows, err := tx.Query("INSERT INTO message (chat_id, by, data) VALUES ($1, $2, $3) RETURNING id;", chatId, message.By, message.Data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var messageId int32
	for messageQueryRows.Next() {
		messageQueryRows.Scan(&messageId)
	}
	tx.Commit()

	return &types.NewChatWithMessageResponse{
		ChatId: chatId,
		MessageId: messageId,
	}, nil
}

func (s *Store) GetAllChats(userId string) (*[]types.Chat, error) {
	var chats []types.Chat
	err := s.db.Select(&chats, "SELECT * FROM chat WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return &chats, nil
}

func (s *Store) GetAllChatMessages(chatId string, userId string) (*[]types.Message, error) {

	var messages []types.Message
	err := s.db.Select(&messages, "SELECT m.id, m.chat_id, m.by, m.data, m.created_at FROM chat c INNER JOIN message m ON c.id = m.chat_id WHERE c.id = $1 AND c.user_id = $2;", chatId, userId)
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("chat not found")
	}

	return &messages, nil
}

func (s *Store) DeleteMessage(userId string, messageId string) error {

	rows, err := s.db.Query("SELECT * FROM message m JOIN chat c ON c.id = m.chat_id WHERE m.id = $1 AND c.user_id = $2;", messageId, userId)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return fmt.Errorf("message not found")
	}

	_, err = s.db.Exec("DELETE FROM message WHERE id = $1", messageId)
	if err != nil {
		return fmt.Errorf("message could not be deleted: %s", err.Error())
	}

	return nil
}

func (s *Store) DeleteChat(userId string, chatId string) error {
	_, err := s.db.Exec("DELETE FROM chat WHERE user_id = $1 AND id = $2", userId, chatId)
	if err != nil {
		return err
	}
	return nil
}
