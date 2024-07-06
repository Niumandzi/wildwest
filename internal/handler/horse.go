package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"wildwest/internal/model/horse"
	"wildwest/internal/service"
	"wildwest/pkg/contextutils"
	"wildwest/pkg/logging"
)

type horseHandler struct {
	service service.HorseService
	logger  logging.Logger
}

func NewHorseHandler(service service.HorseService, logger logging.Logger) HorseHandler {
	return &horseHandler{service: service, logger: logger}
}

// GetHorse retrieves a horse's record by user ID and calculates its speed.
// @Summary Retrieve horse by user ID
// @Description Fetches the horse's data and calculates its speed based on the user ID provided in the path.
// @Tags horse
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} horse.BaseResponse "Returns the horse along with its calculated speed."
// @Failure 400 {string} string "Bad request - user_id is required or invalid."
// @Failure 404 {string} string "Not found - no horse found for the user ID."
// @Router /horse/{user_id} [get]
func (h *horseHandler) GetHorse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	ctx := contextutils.NewContext(r, userID, "GetHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	horseData, err := h.service.GetHorse(ctx, userID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(horseData)
}

// UpgradeHorse upgrades a horse's level for a specified user ID.
// @Summary Upgrade horse level
// @Description Increases the horse's level for the user ID provided in the path.
// @Tags horse
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Horse upgraded successfully with new level information."
// @Failure 400 {string} string "Bad request - user_id is required or invalid."
// @Failure 404 {string} string "Not found - no horse found to upgrade for the user ID."
// @Failure 500 {string} string "Internal server error - error during the upgrade process."
// @Router /horse/upgrade/{user_id} [get]
func (h *horseHandler) UpgradeHorse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	ctx := contextutils.NewContext(r, userID, "UpgradeHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	newLevel, err := h.service.UpgradeHorse(ctx, userID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Horse upgraded successfully",
		"newLevel": newLevel,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GameHorse finishes a horse race and updates the horse's record and earnings based on the distance covered.
// @Summary Finish horse race
// @Description Completes the race for the horse and updates its record and earnings according to the distance covered, as specified in the request body.
// @Tags horse
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param body body horse.GameRequest true "Race completion details including the distance covered."
// @Success 200 {object} horse.GameResponse "Returns the results of the race finish with updated horse data and earnings."
// @Failure 400 {string} string "Bad request - user_id is required or invalid, or the request body is malformed."
// @Failure 404 {string} string "Not found - no horse found to finish race for the user ID."
// @Failure 500 {string} string "Internal server error - error during the race finish process."
// @Router /horse/finish/{user_id} [post]
func (h *horseHandler) GameHorse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var requestData horse.GameRequest
	if err = json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := contextutils.NewContext(r, userID, "GameHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	finishHorseResult, err := h.service.GameHorse(ctx, userID, requestData)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finishHorseResult)
}
