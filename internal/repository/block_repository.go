package repository

import (
	"database/sql"
	"fmt"

	"github.com/dwipurnomo515/yard-planning/internal/model"
)

type BlockRepository struct {
	db *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

// GetByYardAndCode retrieves a block by yard ID and block code
func (r *BlockRepository) GetByYardAndCode(yardID int, code string) (*model.Block, error) {
	query := `
		SELECT id, yard_id, code, name, max_slot, max_row, max_tier, created_at, updated_at
		FROM blocks
		WHERE yard_id = $1 AND code = $2
	`

	var block model.Block
	err := r.db.QueryRow(query, yardID, code).Scan(
		&block.ID,
		&block.YardID,
		&block.Code,
		&block.Name,
		&block.MaxSlot,
		&block.MaxRow,
		&block.MaxTier,
		&block.CreatedAt,
		&block.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("block with code '%s' not found in yard", code)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying block: %w", err)
	}

	return &block, nil
}

// GetByYardID retrieves all blocks for a specific yard
func (r *BlockRepository) GetByYardID(yardID int) ([]model.Block, error) {
	query := `
		SELECT id, yard_id, code, name, max_slot, max_row, max_tier, created_at, updated_at
		FROM blocks
		WHERE yard_id = $1
		ORDER BY code
	`

	rows, err := r.db.Query(query, yardID)
	if err != nil {
		return nil, fmt.Errorf("error querying blocks: %w", err)
	}
	defer rows.Close()

	var blocks []model.Block
	for rows.Next() {
		var block model.Block
		err := rows.Scan(
			&block.ID,
			&block.YardID,
			&block.Code,
			&block.Name,
			&block.MaxSlot,
			&block.MaxRow,
			&block.MaxTier,
			&block.CreatedAt,
			&block.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning block: %w", err)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// GetByID retrieves a block by ID
func (r *BlockRepository) GetByID(id int) (*model.Block, error) {
	query := `
		SELECT id, yard_id, code, name, max_slot, max_row, max_tier, created_at, updated_at
		FROM blocks
		WHERE id = $1
	`

	var block model.Block
	err := r.db.QueryRow(query, id).Scan(
		&block.ID,
		&block.YardID,
		&block.Code,
		&block.Name,
		&block.MaxSlot,
		&block.MaxRow,
		&block.MaxTier,
		&block.CreatedAt,
		&block.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("block with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying block: %w", err)
	}

	return &block, nil
}
