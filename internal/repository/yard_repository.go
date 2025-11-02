package repository

import (
	"database/sql"
	"fmt"

	"github.com/dwipurnomo515/yard-planning/internal/model"
)

type YardRepository struct {
	db *sql.DB
}

func NewYardRepository(db *sql.DB) *YardRepository {
	return &YardRepository{db: db}
}

// GetByCode retrieves a yard by its code
func (r *YardRepository) GetByCode(code string) (*model.Yard, error) {
	query := `
		SELECT id, code, name, description, created_at, updated_at
		FROM yards
		WHERE code = $1
	`

	var yard model.Yard
	err := r.db.QueryRow(query, code).Scan(
		&yard.ID,
		&yard.Code,
		&yard.Name,
		&yard.Description,
		&yard.CreatedAt,
		&yard.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("yard with code '%s' not found", code)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying yard: %w", err)
	}

	return &yard, nil
}

// GetAll retrieves all yards
func (r *YardRepository) GetAll() ([]model.Yard, error) {
	query := `
		SELECT id, code, name, description, created_at, updated_at
		FROM yards
		ORDER BY code
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying yards: %w", err)
	}
	defer rows.Close()

	var yards []model.Yard
	for rows.Next() {
		var yard model.Yard
		err := rows.Scan(
			&yard.ID,
			&yard.Code,
			&yard.Name,
			&yard.Description,
			&yard.CreatedAt,
			&yard.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning yard: %w", err)
		}
		yards = append(yards, yard)
	}

	return yards, nil
}
