package service

import (
	"fmt"

	"github.com/dwipurnomo515/yard-planning/internal/model"
	"github.com/dwipurnomo515/yard-planning/internal/repository"
)

type ContainerService struct {
	yardRepo      *repository.YardRepository
	blockRepo     *repository.BlockRepository
	planRepo      *repository.YardPlanRepository
	containerRepo *repository.ContainerRepository
}

func NewContainerService(
	yardRepo *repository.YardRepository,
	blockRepo *repository.BlockRepository,
	planRepo *repository.YardPlanRepository,
	containerRepo *repository.ContainerRepository,
) *ContainerService {
	return &ContainerService{
		yardRepo:      yardRepo,
		blockRepo:     blockRepo,
		planRepo:      planRepo,
		containerRepo: containerRepo,
	}
}

// GetSuggestion suggests a position for a container based on yard plans
func (s *ContainerService) GetSuggestion(req model.SuggestionRequest) (*model.Position, error) {
	// Validate input
	if err := s.validateContainerSpec(req.ContainerSize, req.ContainerHeight, req.ContainerType); err != nil {
		return nil, err
	}

	// Get yard
	yard, err := s.yardRepo.GetByCode(req.Yard)
	if err != nil {
		return nil, err
	}

	// Get blocks in yard
	blocks, err := s.blockRepo.GetByYardID(yard.ID)
	if err != nil {
		return nil, err
	}

	// Try to find available position in each block
	for _, block := range blocks {
		// Find matching yard plan
		plan, err := s.planRepo.FindMatchingPlan(
			block.ID,
			req.ContainerSize,
			req.ContainerHeight,
			req.ContainerType,
		)
		if err != nil {
			continue // Try next block
		}

		// Find available position in this plan
		position := s.findAvailablePosition(block, *plan)
		if position != nil {
			return position, nil
		}
	}

	return nil, fmt.Errorf("no available position found for container")
}

// PlaceContainer places a container at a specific position
func (s *ContainerService) PlaceContainer(req model.PlacementRequest) error {
	// Validate input
	if req.ContainerNumber == "" {
		return fmt.Errorf("container number is required")
	}

	// Get yard
	yard, err := s.yardRepo.GetByCode(req.Yard)
	if err != nil {
		return err
	}

	// Get block
	block, err := s.blockRepo.GetByYardAndCode(yard.ID, req.Block)
	if err != nil {
		return err
	}

	// Validate position
	if err := s.validatePosition(block, req.Slot, req.Row, req.Tier); err != nil {
		return err
	}

	// Check if container already exists
	existingContainer, _ := s.containerRepo.GetByNumber(req.ContainerNumber)
	if existingContainer != nil {
		return fmt.Errorf("container '%s' already placed in yard", req.ContainerNumber)
	}

	// For now, we'll use default container specs (20ft, 8.6, DRY)
	// In production, you'd want to pass these in the request
	containerSize := 20
	containerHeight := 8.6
	containerType := "DRY"

	// Check if position is available
	occupied, err := s.containerRepo.IsPositionOccupied(block.ID, req.Slot, req.Row, req.Tier, containerSize)
	if err != nil {
		return err
	}
	if occupied {
		return fmt.Errorf("position is already occupied")
	}

	// Check if tier > 1, ensure tier below is occupied
	if req.Tier > 1 {
		occupied, err := s.containerRepo.IsPositionOccupied(block.ID, req.Slot, req.Row, req.Tier-1, containerSize)
		if err != nil {
			return err
		}
		if !occupied {
			return fmt.Errorf("cannot place container at tier %d: tier below is empty", req.Tier)
		}
	}

	// Create container
	container := &model.Container{
		ContainerNumber: req.ContainerNumber,
		YardID:          yard.ID,
		BlockID:         block.ID,
		Slot:            req.Slot,
		Row:             req.Row,
		Tier:            req.Tier,
		ContainerSize:   containerSize,
		ContainerHeight: containerHeight,
		ContainerType:   containerType,
	}

	return s.containerRepo.Create(container)
}

// PickupContainer removes a container from the yard
func (s *ContainerService) PickupContainer(req model.PickupRequest) error {
	// Validate input
	if req.ContainerNumber == "" {
		return fmt.Errorf("container number is required")
	}

	// Get yard
	_, err := s.yardRepo.GetByCode(req.Yard)
	if err != nil {
		return err
	}

	// Get container
	container, err := s.containerRepo.GetByNumber(req.ContainerNumber)
	if err != nil {
		return err
	}

	// Check if container is blocked (has containers on top)
	blocked, err := s.containerRepo.IsContainerBlocked(
		container.BlockID,
		container.Slot,
		container.Row,
		container.Tier,
	)
	if err != nil {
		return err
	}
	if blocked {
		return fmt.Errorf("cannot pickup container: there are containers on top")
	}

	// Delete container
	return s.containerRepo.Delete(req.ContainerNumber)
}

// Helper methods

func (s *ContainerService) validateContainerSpec(size int, height float64, containerType string) error {
	if size != 20 && size != 40 {
		return fmt.Errorf("invalid container size: must be 20 or 40")
	}
	if height != 8.6 && height != 9.6 {
		return fmt.Errorf("invalid container height: must be 8.6 or 9.6")
	}
	validTypes := map[string]bool{"DRY": true, "REEFER": true, "OPEN_TOP": true}
	if !validTypes[containerType] {
		return fmt.Errorf("invalid container type: must be DRY, REEFER, or OPEN_TOP")
	}
	return nil
}

func (s *ContainerService) validatePosition(block *model.Block, slot, row, tier int) error {
	if slot < 1 || slot > block.MaxSlot {
		return fmt.Errorf("invalid slot: must be between 1 and %d", block.MaxSlot)
	}
	if row < 1 || row > block.MaxRow {
		return fmt.Errorf("invalid row: must be between 1 and %d", block.MaxRow)
	}
	if tier < 1 || tier > block.MaxTier {
		return fmt.Errorf("invalid tier: must be between 1 and %d", block.MaxTier)
	}
	return nil
}

func (s *ContainerService) findAvailablePosition(block model.Block, plan model.YardPlan) *model.Position {
	// Get all occupied positions in this plan's area
	occupied, err := s.containerRepo.GetOccupiedPositionsInArea(
		block.ID,
		plan.SlotStart,
		plan.SlotEnd,
		plan.RowStart,
		plan.RowEnd,
	)
	if err != nil {
		return nil
	}

	// Create map of occupied positions
	occupiedMap := make(map[string]bool)
	for _, c := range occupied {
		key := fmt.Sprintf("%d-%d-%d", c.Slot, c.Row, c.Tier)
		occupiedMap[key] = true

		// For 40ft containers, mark second slot as occupied too
		if c.ContainerSize == 40 {
			key2 := fmt.Sprintf("%d-%d-%d", c.Slot+1, c.Row, c.Tier)
			occupiedMap[key2] = true
		}
	}

	// Find first available position (tier 1 first, then stack up)
	for tier := 1; tier <= block.MaxTier; tier++ {
		for slot := plan.SlotStart; slot <= plan.SlotEnd; slot++ {
			for row := plan.RowStart; row <= plan.RowEnd; row++ {
				// For 40ft, need to check if both slots are available
				slotsNeeded := 1
				if plan.ContainerSize == 40 {
					slotsNeeded = 2
					// Make sure we don't exceed slot range
					if slot+1 > plan.SlotEnd {
						continue
					}
				}

				// Check if position is available
				available := true
				for s := 0; s < slotsNeeded; s++ {
					key := fmt.Sprintf("%d-%d-%d", slot+s, row, tier)
					if occupiedMap[key] {
						available = false
						break
					}
				}

				if !available {
					continue
				}

				// For tier > 1, check if tier below is occupied
				if tier > 1 {
					tierBelowOccupied := true
					for s := 0; s < slotsNeeded; s++ {
						keyBelow := fmt.Sprintf("%d-%d-%d", slot+s, row, tier-1)
						if !occupiedMap[keyBelow] {
							tierBelowOccupied = false
							break
						}
					}
					if !tierBelowOccupied {
						continue
					}
				}

				// Found available position
				return &model.Position{
					Block: block.Code,
					Slot:  slot,
					Row:   row,
					Tier:  tier,
				}
			}
		}
	}

	return nil
}
