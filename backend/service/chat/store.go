package chat

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) NewChat(userId string) (int32, error) {

	rows, err := s.db.Query("INSERT INTO chat (user_id) VALUES ($1) RETURNING id;", userId)
	if err != nil {
		return -1, err
	}

	var id int32
	for rows.Next() {
		rows.Scan(&id)
	}
	
	return id, nil
}
