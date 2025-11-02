package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/dwipurnomo515/yard-planning/internal/model"
	"github.com/dwipurnomo515/yard-planning/internal/service"
	"github.com/dwipurnomo515/yard-planning/pkg/response"
	"github.com/dwipurnomo515/yard-planning/pkg/worker"
)

type BulkHandler struct {
	service *service.ContainerService
}

func NewBulkHandler(service *service.ContainerService) *BulkHandler {
	return &BulkHandler{service: service}
}

// BulkSuggestionRequest represents bulk suggestion request
type BulkSuggestionRequest struct {
	Containers []model.SuggestionRequest `json:"containers"`
}

// BulkSuggestionResponse represents bulk suggestion response
type BulkSuggestionResponse struct {
	Results []SuggestionResult `json:"results"`
}

type SuggestionResult struct {
	ContainerNumber   string          `json:"container_number"`
	SuggestedPosition *model.Position `json:"suggested_position,omitempty"`
	Error             string          `json:"error,omitempty"`
}

// HandleBulkSuggestion handles bulk suggestion requests concurrently
func (h *BulkHandler) HandleBulkSuggestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, http.ErrNotSupported)
		return
	}

	var req BulkSuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if len(req.Containers) == 0 {
		response.Error(w, http.StatusBadRequest, http.ErrMissingBoundary)
		return
	}

	// Use worker pool for concurrent processing
	pool := worker.NewPool(5, func(ctx context.Context, job worker.Job) (interface{}, error) {
		suggReq := job.Payload.(model.SuggestionRequest)
		position, err := h.service.GetSuggestion(suggReq)
		if err != nil {
			return nil, err
		}
		return position, nil
	})

	pool.Start()

	// Submit jobs
	go func() {
		for _, container := range req.Containers {
			pool.Submit(worker.Job{
				ID:      container.ContainerNumber,
				Payload: container,
			})
		}
		pool.Stop()
	}()

	// Collect results
	results := make([]SuggestionResult, 0, len(req.Containers))
	for result := range pool.Results() {
		suggReq := result.Job.Payload.(model.SuggestionRequest)
		if result.Err != nil {
			results = append(results, SuggestionResult{
				ContainerNumber: suggReq.ContainerNumber,
				Error:           result.Err.Error(),
			})
		} else {
			position := result.Value.(*model.Position)
			results = append(results, SuggestionResult{
				ContainerNumber:   suggReq.ContainerNumber,
				SuggestedPosition: position,
			})
		}
	}

	resp := BulkSuggestionResponse{Results: results}
	response.Success(w, resp)
}

// BulkPlacementRequest represents bulk placement request
type BulkPlacementRequest struct {
	Containers []model.PlacementRequest `json:"containers"`
}

type PlacementResult struct {
	ContainerNumber string `json:"container_number"`
	Success         bool   `json:"success"`
	Error           string `json:"error,omitempty"`
}

type BulkPlacementResponse struct {
	Results []PlacementResult `json:"results"`
}

// HandleBulkPlacement handles bulk placement with concurrent execution
func (h *BulkHandler) HandleBulkPlacement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, http.ErrNotSupported)
		return
	}

	var req BulkPlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if len(req.Containers) == 0 {
		response.Error(w, http.StatusBadRequest, http.ErrMissingBoundary)
		return
	}

	// Use goroutines with mutex for concurrent placement
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []PlacementResult
	)

	// Limit concurrency to avoid overwhelming the database
	semaphore := make(chan struct{}, 10)

	for _, container := range req.Containers {
		wg.Add(1)
		go func(c model.PlacementRequest) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := h.service.PlaceContainer(c)

			mu.Lock()
			if err != nil {
				results = append(results, PlacementResult{
					ContainerNumber: c.ContainerNumber,
					Success:         false,
					Error:           err.Error(),
				})
			} else {
				results = append(results, PlacementResult{
					ContainerNumber: c.ContainerNumber,
					Success:         true,
				})
			}
			mu.Unlock()
		}(container)
	}

	wg.Wait()

	resp := BulkPlacementResponse{Results: results}
	response.Success(w, resp)
}
