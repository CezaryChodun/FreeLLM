package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cezarychodun/freellms/internal/modules/modelgroups"
)

type ModelsHandler struct {
	repo *modelgroups.ModelGroupRepository
}

func NewModelsHandler(repo *modelgroups.ModelGroupRepository) *ModelsHandler {
	return &ModelsHandler{repo: repo}
}

func (h *ModelsHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listing models")
	groups, err := h.repo.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make([]map[string]any, 0, len(groups)+1)
	data = append(data, map[string]any{
		"id":       "all",
		"object":   "model",
		"created":  time.Now().Unix(),
		"owned_by": "freellm",
	})
	for _, g := range groups {
		data = append(data, map[string]any{
			"id":       g.Name,
			"object":   "model",
			"created":  time.Now().Unix(),
			"owned_by": "freellm",
		})
	}

	resp := map[string]any{
		"object": "list",
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
