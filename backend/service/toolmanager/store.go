package toolmanager

import (
	"database/sql"
	"strings"

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

func (s *Store) GetToolByIdAndUserId(id string, userId string) (*types.Tool, error) {
	dt := new(types.DBTool)
	if err := s.db.Get(dt, "SELECT * FROM tools WHERE id = ? AND user_id = ?;", id, userId); err != nil {
		return nil, err
	}

	deps := []string{}

	if dt.Dependencies != "" {
		deps = strings.Split(dt.Dependencies, ",")
	}

	t := types.Tool{
		ID:           dt.ID,
		Name:         dt.Name,
		Description:  dt.Description,
		Language:     dt.Language,
		Code:         dt.Code,
		Dependencies: deps,
		UserId:       dt.UserId,
	}

	return &t, nil
}

func (s *Store) GetToolsByUserId(userId string) (*[]types.Tool, error) {

	dbTools := []types.DBTool{}
	if err := s.db.Select(&dbTools, "SELECT * FROM tools WHERE user_id = ?;", userId); err != nil {
		return nil, err
	}

	tools := []types.Tool{}
	for _, dt := range dbTools {
		deps := []string{}
		if dt.Dependencies != "" {
			deps = strings.Split(dt.Dependencies, ",")
		}

		tools = append(tools, types.Tool{
			ID:           dt.ID,
			Name:         dt.Name,
			Description:  dt.Description,
			Language:     dt.Language,
			Code:         dt.Code,
			Dependencies: deps,
			UserId:       dt.UserId,
		})
	}

	return &tools, nil
}

func (s *Store) CreateToolByUser(tool *types.CreateTool, userId string) error {
	_, err := s.db.Exec("INSERT INTO tools (name, description, language, code, dependencies, user_id) VALUES (?, ?, ?, ?, ?, ?)",
		tool.Name,
		tool.Description,
		tool.Language,
		tool.Code,
		strings.Join(tool.Dependencies, ","),
		userId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateTool(tool *types.UpdateTool, id string, userId string) error {
	query := "UPDATE tools SET "
	args := []interface{}{}
	fields := []string{}

	if tool.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *tool.Name)
	}
	if tool.Description != nil {
		fields = append(fields, "description = ?")
		args = append(args, *tool.Description)
	}
	if tool.Language != nil {
		fields = append(fields, "language = ?")
		args = append(args, *tool.Language)
	}
	if tool.Code != nil {
		fields = append(fields, "code = ?")
		args = append(args, *tool.Code)
	}
	if tool.Dependencies != nil {
		fields = append(fields, "dependencies = ?")
		args = append(args, strings.Join(*tool.Dependencies, ","))
	}

	if len(fields) == 0 {
		return nil
	}

	query += strings.Join(fields, ", ")
	query += " WHERE id = ? AND user_id = ?"
	args = append(args, id, userId)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) DeleteTool(id string, userId string) error {
	result, err := s.db.Exec("DELETE FROM tools WHERE id = ? AND user_id = ?", id, userId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
