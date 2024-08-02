package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"wildwest/internal/service"
	"wildwest/pkg/logging"
)

type gunfightHandler struct {
	gunfightService   service.GunfightService
	logger            logging.Logger
	upgrader          websocket.Upgrader
	searchConnections sync.Map // Соединения для поиска игры
	gameConnections   sync.Map // Соединения для игры
}

func NewGunfightHandler(gunfightService service.GunfightService, logger logging.Logger) *gunfightHandler {
	return &gunfightHandler{
		gunfightService: gunfightService,
		logger:          logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// FindGunfight initiates a search for an opponent in a gunfight
// @Summary Initiate gunfight search
// @Description Opens a websocket connection and waits to match with an opponent for a gunfight.
// @Tags gunfight
// @Accept json
// @Produce json
// @Param user-id header int true "User ID"
// @Success 200 {object} string "WebSocket connection established, waiting for opponent."
// @Failure 400 {object} string "Could not open websocket connection"
// @Failure 500 {object} string "Internal server error"
// @Router /gunfight/find [get]
func (h *gunfightHandler) FindGunfight(w http.ResponseWriter, r *http.Request) {
	userID, err := h.extractUserID(r)
	if err != nil {
		h.logger.Error("Error extracting user ID: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := h.upgradeConnection(w, r)
	if err != nil {
		h.logger.Error("Error upgrading connection: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	h.searchConnections.Store(userID, conn)
	defer h.searchConnections.Delete(userID)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				h.logger.Error("Error reading from websocket: ", err)
				conn.Close()
				cancel()
				break
			}
		}
	}()

	result, err := h.gunfightService.FindGunfight(ctx, userID)
	if err != nil {
		h.logger.Error("Error finding gunfight: ", err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	if result.OpponentID != 0 {
		message := fmt.Sprintf("Gunfight ID: %s", result.Message)
		h.sendMessageAndClose(conn, message, userID)

		if opponentConn, ok := h.searchConnections.Load(result.OpponentID); ok {
			h.sendMessageAndClose(opponentConn.(*websocket.Conn), message, result.OpponentID)
		}
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte(result.Message))
	}

	if ctx.Err() != nil {
		h.gunfightService.RemovePlayerFromQueue(context.Background(), userID)
	}
}

func (h *gunfightHandler) sendMessageAndClose(conn *websocket.Conn, message string, userID int) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	conn.Close()
	h.searchConnections.Delete(userID)
}

func (h *gunfightHandler) extractUserID(r *http.Request) (int, error) {
	userData, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("user data is required")
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user ID is required and must be a number")
	}

	return int(userID), nil
}

func (h *gunfightHandler) upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		h.logger.Error("Could not open websocket connection: ", err)
		return nil, err
	}
	return conn, nil
}
