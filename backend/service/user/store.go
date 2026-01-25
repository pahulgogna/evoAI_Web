package user

import (
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

func (s *Store) GetUserById(id string) (*types.User, error) {
	user := new(types.User)

	if err := s.db.Get(user, "SELECT * FROM users WHERE id = $1;", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	user := new(types.User)

	if err := s.db.Get(user, "SELECT * FROM users WHERE email = $1;", email); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) CreateUser(user *types.User) error {
	_, err := s.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	return err
}

func (s *Store) RemoveUser(id string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1;", id)
	return err
}
