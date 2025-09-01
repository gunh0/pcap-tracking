package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"pcap-tracking/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

// Client represents a websocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	running    bool
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	h.mutex.Lock()
	h.running = true
	h.mutex.Unlock()

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			logrus.WithField("client_count", len(h.clients)).Info("Client connected")

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			logrus.WithField("client_count", len(h.clients)).Info("Client disconnected")

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mutex.RUnlock()
		}

		h.mutex.RLock()
		running := h.running
		h.mutex.RUnlock()

		if !running {
			break
		}
	}
}

// Stop stops the hub
func (h *Hub) Stop() {
	h.mutex.Lock()
	h.running = false

	// Close all client connections
	for client := range h.clients {
		delete(h.clients, client)
		close(client.send)
		client.conn.Close()
	}
	h.mutex.Unlock()
}

// BroadcastStats broadcasts packet statistics to all connected clients
func (h *Hub) BroadcastStats(stats *models.SystemStats) {
	message := models.WebSocketMessage{
		Type: "stats",
		Data: stats,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal stats message")
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel full, dropping stats message")
	}
}

// BroadcastPacket broadcasts individual packet information to all connected clients
func (h *Hub) BroadcastPacket(packet *models.PacketInfo) {
	message := models.WebSocketMessage{
		Type: "packet",
		Data: packet,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal packet message")
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel full, dropping packet message")
	}
}

// HandleWebSocket handles websocket connections
func (h *Hub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade websocket connection")
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

	// Start goroutines for this client
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("Websocket error")
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logrus.WithError(err).Error("Failed to write message")
				return
			}
		}
	}
}
