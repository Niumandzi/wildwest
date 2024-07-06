package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"wildwest/internal/service"
	"wildwest/pkg/contextutils"
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
	conn, err := h.upgradeConnection(w, r)
	if err != nil {
		return
	}
	userID, err := h.readUserID(conn)
	if err != nil {
		return
	}

	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	notifyChan := newSafeChannel()
	defer notifyChan.close()

	ctx := contextutils.NewContext(r, userID, "FindGunfight")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	h.handleConnection(ctx, conn, userID, notifyChan)
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

func (h *gunfightHandler) readUserID(conn *websocket.Conn) (int, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to read message"))
		h.logger.Error("Failed to read message: ", err)
		conn.Close()
		return 0, err
	}
	userID, err := strconv.Atoi(string(message))
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid userID format"))
		h.logger.Error("Invalid userID format: ", err)
		conn.Close()
		return 0, err
	}
	return userID, nil
}

func (h *gunfightHandler) handleConnection(ctx context.Context, conn *websocket.Conn, userID int, notifyChan *safeChannel) {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go h.waitForOpponent(cancelCtx, conn, userID, notifyChan)
	go h.monitorConnection(conn, cancel)

	select {
	case opponentID := <-notifyChan.C:
		h.notifyPlayers(conn, userID, opponentID)
	case <-cancelCtx.Done():
		h.logger.Info("Connection closed by client or error occurred")
		h.cleanupAfterDisconnect(userID)
	case <-time.After(1 * time.Minute):
		h.logger.Info("No opponent found within the time limit")
		conn.WriteMessage(websocket.TextMessage, []byte("No opponent found within the time limit"))
		conn.Close()
	}
}

func (h *gunfightHandler) waitForOpponent(ctx context.Context, conn *websocket.Conn, userID int, notifyChan *safeChannel) {
	opponentID, err := h.gunfightService.FindGunfight(ctx, userID, notifyChan.C)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		h.logger.Error("Error finding gunfight: ", err)
		conn.Close()
		return
	}
	if ctx.Err() != nil {
		return
	}
	if err := notifyChan.send(opponentID); err != nil {
		h.logger.Error("Failed to send opponent ID: ", err)
	}
}

func (h *gunfightHandler) monitorConnection(conn *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.logger.Error("Connection error or closed by client: ", err)
			return
		}
	}
}

func (h *gunfightHandler) notifyPlayers(conn *websocket.Conn, userID, opponentID int) {
	message := fmt.Sprintf("Opponent found: %d", opponentID)
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	if opponentConn, ok := h.connections.Load(opponentID); ok {
		opponentConn.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Opponent found: %d", userID)))
		opponentConn.(*websocket.Conn).Close()
	}
	conn.Close()
}

func (h *gunfightHandler) cleanupAfterDisconnect(userID int) {
	h.gunfightService.RemovePlayerFromQueue(context.Background(), userID)
	if conn, ok := h.connections.Load(userID); ok {
		conn.(*websocket.Conn).Close()
	}
}

type safeChannel struct {
	C      chan int
	closed bool
	mutex  sync.Mutex
}

func newSafeChannel() *safeChannel {
	return &safeChannel{
		C:      make(chan int),
		closed: false,
	}
}

func (sc *safeChannel) send(value int) error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if sc.closed {
		return fmt.Errorf("channel is closed")
	}
	sc.C <- value
	return nil
}

func (sc *safeChannel) close() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if !sc.closed {
		close(sc.C)
		sc.closed = true
	}
}
