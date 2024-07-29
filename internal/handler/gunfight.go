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
	gunfightService service.GunfightService
	logger          logging.Logger
	upgrader        websocket.Upgrader
	connections     sync.Map
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := h.upgradeConnection(w, r)
	if err != nil {
		return
	}
	defer conn.Close()

	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Мониторинг закрытия соединения
	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			cancel() // Отменяем контекст, если соединение закрыто
		}
	}()

	defer h.handleUserDisconnect(ctx, userID)

	opponentID, err := h.gunfightService.FindGunfight(ctx, userID)
	if err != nil {
		if err.Error() == "no opponent found within the time limit" {
			h.handleTimeout(conn)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return
	}

	h.handleOpponentsFound(conn, userID, opponentID)
}

func (h *gunfightHandler) handleUserDisconnect(ctx context.Context, userID int) {
	if ctx.Err() != nil {
		h.gunfightService.RemovePlayerFromQueue(context.Background(), userID)
	}
}

func (h *gunfightHandler) handleTimeout(conn *websocket.Conn) {
	message := "No opponent found within the time limit"
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	conn.Close()
}

func (h *gunfightHandler) handleOpponentsFound(conn *websocket.Conn, userID, opponentID int) {
	// Отправка сообщения текущему пользователю и закрытие соединения
	message := fmt.Sprintf("Opponent found: %d", opponentID)
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	conn.Close()

	// Отправка сообщения оппоненту и закрытие его соединения
	if opponentConn, ok := h.connections.Load(opponentID); ok {
		opponentConn.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Opponent found: %d", userID)))
		opponentConn.(*websocket.Conn).Close()
		h.connections.Delete(opponentID)
	}
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
