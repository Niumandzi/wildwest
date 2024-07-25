package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"wildwest/internal/service"
	"wildwest/pkg/contextutils"
	"wildwest/pkg/logging"
)

type moneyHandler struct {
	service service.MoneyService
	logger  logging.Logger
}

func NewMoneyHandler(service service.MoneyService, logger logging.Logger) MoneyHandler {
	return &moneyHandler{service: service, logger: logger}
}

// GetMoney retrieves a user's money record by user ID.
// @Summary Retrieve money record
// @Description Fetches the money details associated with the provided user ID.
// @Tags money
// @Accept json
// @Produce json
// @Param X-User-Data header string true "User data in encoded format containing user ID and other necessary information"
// @Success 200 {object} money.BaseResponse "Returns the money details for the specified user ID."
// @Failure 400 {string} string "Bad request - user_id is required or invalid."
// @Failure 404 {string} string "Not found - no money record found for the user ID."
// @Router /money [get]
func (h *moneyHandler) GetMoney(w http.ResponseWriter, r *http.Request) {
	userData, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		http.Error(w, "User data is required", http.StatusBadRequest)
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		http.Error(w, "User ID is required and must be a number", http.StatusBadRequest)
		return
	}

	ctx := contextutils.NewContext(r, int(userID), "UpgradeHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	money, err := h.service.GetMoney(ctx, int(userID))
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(money)
}
