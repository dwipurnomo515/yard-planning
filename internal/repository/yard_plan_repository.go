package repository

import (
	"database/sql"
	"fmt"

	"github.com/dwipurnomo515/yard-planning/internal/model"
)

type YardPlanRepository struct {
	db *sql.DB
}

func NewYardPlanRepository(db *sql.DB) *YardPlanRepository {
	return &YardPlanRepository{db: db}
}

// FindMatchingPlan finds a yard plan that matches the container specifications
func (r *YardPlanRepository) FindMatchingPlan(blockID int, size int, height float64, containerType string) (*model.YardPlan, error) {
	query := `
		SELECT id, block_id, slot_start, slot_end, row_start, row_end,
		       container_size, container_height, container_type, stacking_priority,
		       created_at, updated_at
		FROM yard_plans
		WHERE block_id = $1
		  AND container_size = $2
		  AND container_height = $3
		  AND container_type = $4
		LIMIT 1
	`

	var plan model.YardPlan
	err := r.db.QueryRow(query, blockID, size, height, containerType).Scan(
		&plan.ID,
		&plan.BlockID,
		&plan.SlotStart,
		&plan.SlotEnd,
		&plan.RowStart,
		&plan.RowEnd,
		&plan.ContainerSize,
		&plan.ContainerHeight,
		&plan.ContainerType,
		&plan.StackingPriority,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no yard plan found for container size=%d, height=%.1f, type=%s", size, height, containerType)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying yard plan: %w", err)
	}

	return &plan, nil
}

// GetByBlockID retrieves all yard plans for a specific block
func (r *YardPlanRepository) GetByBlockID(blockID int) ([]model.YardPlan, error) {
	query := `
		SELECT id, block_id, slot_start, slot_end, row_start, row_end,
		       container_size, container_height, container_type, stacking_priority,
		       created_at, updated_at
		FROM yard_plans
		WHERE block_id = $1
		ORDER BY slot_start, row_start
	`

	rows, err := r.db.Query(query, blockID)
	if err != nil {
		return nil, fmt.Errorf("error querying yard plans: %w", err)
	}
	defer rows.Close()

	var plans []model.YardPlan
	for rows.Next() {
		var plan model.YardPlan
		err := rows.Scan(
			&plan.ID,
			&plan.BlockID,
			&plan.SlotStart,
			&plan.SlotEnd,
			&plan.RowStart,
			&plan.RowEnd,
			&plan.ContainerSize,
			&plan.ContainerHeight,
			&plan.ContainerType,
			&plan.StackingPriority,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning yard plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// Create creates a new yard plan
func (r *YardPlanRepository) Create(plan *model.YardPlan) error {
	query := `
		INSERT INTO yard_plans (
			block_id, slot_start, slot_end, row_start, row_end,
			container_size, container_height, container_type, stacking_priority
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		plan.BlockID,
		plan.SlotStart,
		plan.SlotEnd,
		plan.RowStart,
		plan.RowEnd,
		plan.ContainerSize,
		plan.ContainerHeight,
		plan.ContainerType,
		plan.StackingPriority,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating yard plan: %w", err)
	}

	return nil
}
