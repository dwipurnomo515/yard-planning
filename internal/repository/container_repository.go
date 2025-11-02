package repository

import (
	"database/sql"
	"fmt"

	"github.com/dwipurnomo515/yard-planning/internal/model"
)

type ContainerRepository struct {
	db *sql.DB
}

func NewContainerRepository(db *sql.DB) *ContainerRepository {
	return &ContainerRepository{db: db}
}

// Create inserts a new container into the database
func (r *ContainerRepository) Create(container *model.Container) error {
	query := `
		INSERT INTO containers (
			container_number, yard_id, block_id, slot, row, tier,
			container_size, container_height, container_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, placed_at
	`

	err := r.db.QueryRow(
		query,
		container.ContainerNumber,
		container.YardID,
		container.BlockID,
		container.Slot,
		container.Row,
		container.Tier,
		container.ContainerSize,
		container.ContainerHeight,
		container.ContainerType,
	).Scan(&container.ID, &container.PlacedAt)

	if err != nil {
		return fmt.Errorf("error creating container: %w", err)
	}

	return nil
}

// GetByNumber retrieves a container by its number
func (r *ContainerRepository) GetByNumber(containerNumber string) (*model.Container, error) {
	query := `
		SELECT id, container_number, yard_id, block_id, slot, row, tier,
		       container_size, container_height, container_type, placed_at
		FROM containers
		WHERE container_number = $1
	`

	var container model.Container
	err := r.db.QueryRow(query, containerNumber).Scan(
		&container.ID,
		&container.ContainerNumber,
		&container.YardID,
		&container.BlockID,
		&container.Slot,
		&container.Row,
		&container.Tier,
		&container.ContainerSize,
		&container.ContainerHeight,
		&container.ContainerType,
		&container.PlacedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("container '%s' not found", containerNumber)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying container: %w", err)
	}

	return &container, nil
}

// Delete removes a container from the database
func (r *ContainerRepository) Delete(containerNumber string) error {
	query := `DELETE FROM containers WHERE container_number = $1`

	result, err := r.db.Exec(query, containerNumber)
	if err != nil {
		return fmt.Errorf("error deleting container: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("container '%s' not found", containerNumber)
	}

	return nil
}

// IsPositionOccupied checks if a specific position is occupied
// For 40ft containers, checks both slots
func (r *ContainerRepository) IsPositionOccupied(blockID, slot, row, tier int, containerSize int) (bool, error) {
	var query string
	var args []interface{}

	if containerSize == 40 {
		// Check both slot and slot+1 for 40ft container
		query = `
			SELECT COUNT(*) > 0
			FROM containers
			WHERE block_id = $1
			  AND row = $2
			  AND tier = $3
			  AND (slot = $4 OR slot = $5)
		`
		args = []interface{}{blockID, row, tier, slot, slot + 1}
	} else {
		// Check single slot for 20ft container
		query = `
			SELECT COUNT(*) > 0
			FROM containers
			WHERE block_id = $1
			  AND slot = $2
			  AND row = $3
			  AND tier = $4
		`
		args = []interface{}{blockID, slot, row, tier}
	}

	var occupied bool
	err := r.db.QueryRow(query, args...).Scan(&occupied)
	if err != nil {
		return false, fmt.Errorf("error checking position: %w", err)
	}

	return occupied, nil
}

// GetOccupiedPositionsInArea retrieves all occupied positions within a specific area
func (r *ContainerRepository) GetOccupiedPositionsInArea(blockID, slotStart, slotEnd, rowStart, rowEnd int) ([]model.Container, error) {
	query := `
		SELECT id, container_number, yard_id, block_id, slot, row, tier,
		       container_size, container_height, container_type, placed_at
		FROM containers
		WHERE block_id = $1
		  AND slot >= $2
		  AND slot <= $3
		  AND row >= $4
		  AND row <= $5
		ORDER BY slot, row, tier
	`

	rows, err := r.db.Query(query, blockID, slotStart, slotEnd, rowStart, rowEnd)
	if err != nil {
		return nil, fmt.Errorf("error querying containers: %w", err)
	}
	defer rows.Close()

	var containers []model.Container
	for rows.Next() {
		var container model.Container
		err := rows.Scan(
			&container.ID,
			&container.ContainerNumber,
			&container.YardID,
			&container.BlockID,
			&container.Slot,
			&container.Row,
			&container.Tier,
			&container.ContainerSize,
			&container.ContainerHeight,
			&container.ContainerType,
			&container.PlacedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning container: %w", err)
		}
		containers = append(containers, container)
	}

	return containers, nil
}

// IsContainerBlocked checks if there's a container above the given position
func (r *ContainerRepository) IsContainerBlocked(blockID, slot, row, tier int) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM containers
		WHERE block_id = $1
		  AND slot = $2
		  AND row = $3
		  AND tier > $4
	`

	var blocked bool
	err := r.db.QueryRow(query, blockID, slot, row, tier).Scan(&blocked)
	if err != nil {
		return false, fmt.Errorf("error checking if blocked: %w", err)
	}

	return blocked, nil
}

// GetAll retrieves all containers
func (r *ContainerRepository) GetAll() ([]model.Container, error) {
	query := `
		SELECT id, container_number, yard_id, block_id, slot, row, tier,
		       container_size, container_height, container_type, placed_at
		FROM containers
		ORDER BY placed_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying containers: %w", err)
	}
	defer rows.Close()

	var containers []model.Container
	for rows.Next() {
		var container model.Container
		err := rows.Scan(
			&container.ID,
			&container.ContainerNumber,
			&container.YardID,
			&container.BlockID,
			&container.Slot,
			&container.Row,
			&container.Tier,
			&container.ContainerSize,
			&container.ContainerHeight,
			&container.ContainerType,
			&container.PlacedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning container: %w", err)
		}
		containers = append(containers, container)
	}

	return containers, nil
}

// GetByBlock retrieves all containers in a specific block
func (r *ContainerRepository) GetByBlock(blockID int) ([]model.Container, error) {
	query := `
		SELECT id, container_number, yard_id, block_id, slot, row, tier,
		       container_size, container_height, container_type, placed_at
		FROM containers
		WHERE block_id = $1
		ORDER BY slot, row, tier
	`

	rows, err := r.db.Query(query, blockID)
	if err != nil {
		return nil, fmt.Errorf("error querying containers: %w", err)
	}
	defer rows.Close()

	var containers []model.Container
	for rows.Next() {
		var container model.Container
		err := rows.Scan(
			&container.ID,
			&container.ContainerNumber,
			&container.YardID,
			&container.BlockID,
			&container.Slot,
			&container.Row,
			&container.Tier,
			&container.ContainerSize,
			&container.ContainerHeight,
			&container.ContainerType,
			&container.PlacedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning container: %w", err)
		}
		containers = append(containers, container)
	}

	return containers, nil
}
