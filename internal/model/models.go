package model

import "time"

// Yard represents a container yard
type Yard struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Block represents a storage block in a yard
type Block struct {
	ID        int       `json:"id"`
	YardID    int       `json:"yard_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	MaxSlot   int       `json:"max_slot"`
	MaxRow    int       `json:"max_row"`
	MaxTier   int       `json:"max_tier"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// YardPlan represents a planning configuration for specific area
type YardPlan struct {
	ID               int       `json:"id"`
	BlockID          int       `json:"block_id"`
	SlotStart        int       `json:"slot_start"`
	SlotEnd          int       `json:"slot_end"`
	RowStart         int       `json:"row_start"`
	RowEnd           int       `json:"row_end"`
	ContainerSize    int       `json:"container_size"`
	ContainerHeight  float64   `json:"container_height"`
	ContainerType    string    `json:"container_type"`
	StackingPriority string    `json:"stacking_priority"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Container represents a physical container in the yard
type Container struct {
	ID              int       `json:"id"`
	ContainerNumber string    `json:"container_number"`
	YardID          int       `json:"yard_id"`
	BlockID         int       `json:"block_id"`
	Slot            int       `json:"slot"`
	Row             int       `json:"row"`
	Tier            int       `json:"tier"`
	ContainerSize   int       `json:"container_size"`
	ContainerHeight float64   `json:"container_height"`
	ContainerType   string    `json:"container_type"`
	PlacedAt        time.Time `json:"placed_at"`
}

// Position represents a container position
type Position struct {
	Block string `json:"block"`
	Slot  int    `json:"slot"`
	Row   int    `json:"row"`
	Tier  int    `json:"tier"`
}

// Request/Response DTOs
type SuggestionRequest struct {
	Yard            string  `json:"yard"`
	ContainerNumber string  `json:"container_number"`
	ContainerSize   int     `json:"container_size"`
	ContainerHeight float64 `json:"container_height"`
	ContainerType   string  `json:"container_type"`
}

type SuggestionResponse struct {
	SuggestedPosition Position `json:"suggested_position"`
}

type PlacementRequest struct {
	Yard            string `json:"yard"`
	ContainerNumber string `json:"container_number"`
	Block           string `json:"block"`
	Slot            int    `json:"slot"`
	Row             int    `json:"row"`
	Tier            int    `json:"tier"`
}

type PlacementResponse struct {
	Message string `json:"message"`
}

type PickupRequest struct {
	Yard            string `json:"yard"`
	ContainerNumber string `json:"container_number"`
}

type PickupResponse struct {
	Message string `json:"message"`
}
