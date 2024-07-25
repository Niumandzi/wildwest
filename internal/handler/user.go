package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"wildwest/internal/model/user"
	"wildwest/internal/service"
	"wildwest/pkg/contextutils"
	"wildwest/pkg/logging"
)

type userHandler struct {
	service service.UserService
	logger  logging.Logger
}

func NewUserHandler(service service.UserService, logger logging.Logger) UserHandler {
	return &userHandler{service: service, logger: logger}
}

// CheckUser check user and his actual data.
// @Summary Check user and his actual data
// @Description check user and his actual data.
// @Tags user
// @Accept json
// @Produce json
// @Param X-User-Data header string true "User data in encoded format containing user ID and other necessary information"
// @Success 200 {object} user.BaseResponse "Returns the user details."
// @Failure 400 {string} string "Bad request - invalid request body."
// @Failure 500 {string} string "Internal server error - failed to create or update user."
// @Router /user [get]
func (h *userHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
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

	firstName, ok := userData["first_name"].(string)
	if !ok {
		http.Error(w, "User ID is required and must be a number", http.StatusBadRequest)
		return
	}

	lastName, ok := userData["last_name"].(string)
	if !ok {
		http.Error(w, "User ID is required and must be a number", http.StatusBadRequest)
		return
	}

	username, ok := userData["username"].(string)
	if !ok {
		http.Error(w, "User ID is required and must be a number", http.StatusBadRequest)
		return
	}

	userRequest := user.BaseRequest{
		ID:        int(userID),
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}

	ctx := contextutils.NewContext(r, int(userID), "UpgradeHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := h.service.CreateOrUpdateUser(ctx, userRequest)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to create or update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
