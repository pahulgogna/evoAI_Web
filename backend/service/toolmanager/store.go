package toolmanager

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

func (s *Store) GetToolByName(name string) (*types.Tool, error) {

	t := new(types.Tool)

	if err := s.db.Get(t, "SELECT * FROM tools WHERE name = $1;", name); err != nil {
		return nil, err
	}

	return t, nil
}
