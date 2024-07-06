package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
// @Param user_id path int true "User ID" "The unique identifier of the user to retrieve money details for."
// @Success 200 {object} money.BaseResponse "Returns the money details for the specified user ID."
// @Failure 400 {string} string "Bad request - user_id is required or invalid."
// @Failure 404 {string} string "Not found - no money record found for the user ID."
// @Router /money/{user_id} [get]
func (h *moneyHandler) GetMoney(w http.ResponseWriter, r *http.Request) {
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

	money, err := h.service.GetMoney(ctx, userID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(money)
}
