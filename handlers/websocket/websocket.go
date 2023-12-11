package websocket

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/raihan2bd/vidverse/config"
)

type Repo struct {
	App     *config.Application
	Clients *Clients
}

var Methods *Repo

func NewAPP(a *config.Application) *Repo {
	return &Repo{
		App:     a,
		Clients: &Clients{m: map[uint]*websocket.Conn{}},
	}
}

func NewSocket(m *Repo) {
	Methods = m
}

type Clients struct {
	sync.RWMutex
	m map[uint]*websocket.Conn
}

func (c *Clients) Add(userID uint, conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()
	c.m[userID] = conn
}

func (c *Clients) Remove(userID uint) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, userID)
}

func (c *Clients) Get(userID uint) *websocket.Conn {
	c.RLock()
	defer c.RUnlock()
	return c.m[userID]
}

// handle websocket request
func (m *Repo) WSHandler(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from context or request
	userID := uint(1) // Extract user ID using appropriate method

	// Add client to map
	m.Clients.Add(userID, conn)

	// Send connection event message
	m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "connect"}

	go func() {
		defer func() {
			conn.Close()
			m.Clients.Remove(userID)
			m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "disconnect"}
		}()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				// Handle error and send disconnect event
				m.MessageCh <- &MessageEvent{userID: userID, conn: conn, action: "disconnect", err: err}
				break
			}

			// Handle message based on type
			// ...

			// Broadcast message to other users
			// ...

			// Send message event
			m.MessageCh <- &MessageEvent{userID: userID, conn: conn, action: "message", data: message}
		}
	}()
}

func (m *Repo) HandleMessages() {
	for {
		select {
		case message := <-m.MessageCh:
			switch message.action {
			case "connect":
				// Handle connection establishment
				// ...

			case "message":
				// Handle received message
				// ...

			case "disconnect":
				// Handle client disconnect
				// ...
			}
		default:
			// Optionally: add code to perform background tasks
		}
	}
}
