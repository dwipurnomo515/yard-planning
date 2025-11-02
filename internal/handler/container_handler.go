package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dwipurnomo515/yard-planning/internal/model"
	"github.com/dwipurnomo515/yard-planning/internal/service"
	"github.com/dwipurnomo515/yard-planning/pkg/response"
)

type ContainerHandler struct {
	service *service.ContainerService
}

func NewContainerHandler(service *service.ContainerService) *ContainerHandler {
	return &ContainerHandler{service: service}
}

// HandleSuggestion handles POST /suggestion
func (h *ContainerHandler) HandleSuggestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed,
			http.ErrNotSupported)
		return
	}

	var req model.SuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Validate required fields
	if req.Yard == "" {
		response.Error(w, http.StatusBadRequest,
			http.ErrMissingBoundary)
		return
	}
	if req.ContainerNumber == "" {
		response.Error(w, http.StatusBadRequest,
			http.ErrMissingBoundary)
		return
	}

	position, err := h.service.GetSuggestion(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	resp := model.SuggestionResponse{
		SuggestedPosition: *position,
	}

	response.Success(w, resp)
}

// HandlePlacement handles POST /placement
func (h *ContainerHandler) HandlePlacement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed,
			http.ErrNotSupported)
		return
	}

	var req model.PlacementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Validate required fields
	if req.Yard == "" || req.ContainerNumber == "" || req.Block == "" {
		response.Error(w, http.StatusBadRequest,
			http.ErrMissingBoundary)
		return
	}
	if req.Slot < 1 || req.Row < 1 || req.Tier < 1 {
		response.Error(w, http.StatusBadRequest,
			http.ErrMissingBoundary)
		return
	}

	err := h.service.PlaceContainer(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	resp := model.PlacementResponse{
		Message: "Success",
	}

	response.Success(w, resp)
}

// HandlePickup handles POST /pickup
func (h *ContainerHandler) HandlePickup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed,
			http.ErrNotSupported)
		return
	}

	var req model.PickupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Validate required fields
	if req.Yard == "" || req.ContainerNumber == "" {
		response.Error(w, http.StatusBadRequest,
			http.ErrMissingBoundary)
		return
	}

	err := h.service.PickupContainer(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	resp := model.PickupResponse{
		Message: "Success",
	}

	response.Success(w, resp)
}
