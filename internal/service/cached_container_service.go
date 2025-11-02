package service

import (
	"fmt"
	"time"

	"github.com/dwipurnomo515/yard-planning/internal/model"
	"github.com/dwipurnomo515/yard-planning/internal/repository"
	"github.com/dwipurnomo515/yard-planning/pkg/cache"
)

type CachedContainerService struct {
	ContainerService
	cache *cache.RedisClient
}

func NewCachedContainerService(
	yardRepo *repository.YardRepository,
	blockRepo *repository.BlockRepository,
	planRepo *repository.YardPlanRepository,
	containerRepo *repository.ContainerRepository,
	redisClient *cache.RedisClient,
) *CachedContainerService {
	return &CachedContainerService{
		ContainerService: ContainerService{
			yardRepo:      yardRepo,
			blockRepo:     blockRepo,
			planRepo:      planRepo,
			containerRepo: containerRepo,
		},
		cache: redisClient,
	}
}

// GetSuggestion with caching
func (s *CachedContainerService) GetSuggestion(req model.SuggestionRequest) (*model.Position, error) {
	// Validate input
	if err := s.validateContainerSpec(req.ContainerSize, req.ContainerHeight, req.ContainerType); err != nil {
		return nil, err
	}

	// Try to get from cache
	cacheKey := fmt.Sprintf("suggestion:%s:%d:%.1f:%s",
		req.Yard, req.ContainerSize, req.ContainerHeight, req.ContainerType)

	var cachedPosition model.Position
	if err := s.cache.Get(cacheKey, &cachedPosition); err == nil {
		// Verify position is still available
		yard, _ := s.yardRepo.GetByCode(req.Yard)
		if yard != nil {
			block, _ := s.blockRepo.GetByYardAndCode(yard.ID, cachedPosition.Block)
			if block != nil {
				occupied, _ := s.containerRepo.IsPositionOccupied(
					block.ID,
					cachedPosition.Slot,
					cachedPosition.Row,
					cachedPosition.Tier,
					req.ContainerSize,
				)
				if !occupied {
					return &cachedPosition, nil
				}
			}
		}
		// Cache invalid, delete it
		s.cache.Delete(cacheKey)
	}

	// Get fresh suggestion
	position, err := s.ContainerService.GetSuggestion(req)
	if err != nil {
		return nil, err
	}

	// Cache the result for 5 minutes
	s.cache.Set(cacheKey, position, 5*time.Minute)

	return position, nil
}

// PlaceContainer with cache invalidation
func (s *CachedContainerService) PlaceContainer(req model.PlacementRequest) error {
	err := s.ContainerService.PlaceContainer(req)
	if err != nil {
		return err
	}

	// Invalidate related caches
	pattern := fmt.Sprintf("suggestion:%s:*", req.Yard)
	s.cache.DeletePattern(pattern)

	// Cache the container position
	cacheKey := fmt.Sprintf("container:%s", req.ContainerNumber)
	containerInfo := map[string]interface{}{
		"yard":  req.Yard,
		"block": req.Block,
		"slot":  req.Slot,
		"row":   req.Row,
		"tier":  req.Tier,
	}
	s.cache.Set(cacheKey, containerInfo, 24*time.Hour)

	return nil
}

// PickupContainer with cache invalidation
func (s *CachedContainerService) PickupContainer(req model.PickupRequest) error {
	err := s.ContainerService.PickupContainer(req)
	if err != nil {
		return err
	}

	// Invalidate caches
	pattern := fmt.Sprintf("suggestion:%s:*", req.Yard)
	s.cache.DeletePattern(pattern)

	// Remove container cache
	cacheKey := fmt.Sprintf("container:%s", req.ContainerNumber)
	s.cache.Delete(cacheKey)

	return nil
}
