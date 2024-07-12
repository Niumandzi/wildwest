package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"sort"
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

// Authenticate authenticates a user via Telegram Web App and issues a JWT.
// @Summary Authenticate user via Telegram
// @Description Authenticates a user by verifying the Telegram Web App data and issues a JWT if successful.
// @Tags user
// @Accept json
// @Produce json
// @Param body body user.BaseRequest true "User data required for authentication."
// @Param X-User-Data header string true "User data in encoded format containing user ID and other necessary information"
// @Success 200 {object} map[string]string "JWT token issued successfully."
// @Failure 400 {string} string "Bad request - invalid request body."
// @Failure 401 {string} string "Unauthorized - invalid data signature."
// @Failure 500 {string} string "Internal server error - failed to create or update user, or token generation failed."
// @Router /auth/telegram [post]
func (h *userHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var userRequest user.BaseRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка подписи (реализуйте эту функцию)
	if !checkTelegramAuth(userRequest, "YOUR_BOT_TOKEN") {
		http.Error(w, "Invalid data signature", http.StatusUnauthorized)
		return
	}

	ctx := contextutils.NewContext(r, 0, "UpgradeHorse")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userData, err := h.service.RegisterUser(ctx, userRequest)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to create or update user", http.StatusInternalServerError)
		return
	}

	// Создание JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userData.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Отправка токена пользователю
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func checkTelegramAuth(request user.BaseRequest, botToken string) bool {
	// Collecting data from the request
	data := map[string]string{
		"id":         fmt.Sprintf("%d", request.ID),
		"first_name": request.FirstName,
		"last_name":  request.LastName,
	}

	// Creating the data check string
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var checkString string
	for _, key := range keys {
		if data[key] != "" {
			checkString += fmt.Sprintf("%s=%s\n", key, data[key])
		}
	}

	// Hashing the data check string with the bot token
	hash := hmac.New(sha256.New, []byte(botToken))
	hash.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(hash.Sum(nil))

	// Comparing the computed hash with the provided hash
	return hmac.Equal([]byte(expectedHash), []byte(request.Hash))
}
