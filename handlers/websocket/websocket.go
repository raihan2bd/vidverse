package websocket

import (
	"fmt"
	"net/http"
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

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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

type WsPayload struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// handle websocket request
func (m *Repo) WSHandler(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	fmt.Println("WSHandler")
	conn, err := upgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from context or request
	userID := uint(1)

	// Add client to map
	m.Clients.Add(userID, conn)

	// Send connection event message
	m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "message", Data: "connected"}
	// conn.WriteJSON(WsPayload{Action: "connect", Data: "connected"})

	go func() {
		defer func() {
			conn.Close()
			m.Clients.Remove(userID)
			m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "disconnect"}
		}()

		var payload WsPayload

		for {
			err := conn.ReadJSON(&payload)
			if err != nil {
				// Handle error and send disconnect event
				m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "disconnect"}
				break
			}

			fmt.Println("payload", payload)

			// Send message event
			m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Data: payload, Action: payload.Action}
		}
	}()
}

func (m *Repo) HandleMessages() {
	for {
		event := <-m.App.NotificationChan
		switch event.Action {
		case "connect":
			// get clients conn from map using broadcasterID
			conn := m.Clients.Get(event.BroadcasterID)
			conn.WriteJSON(WsPayload{Action: "connect", Data: "connected"})

		case "a_new_notification":
			conn := m.Clients.Get(event.BroadcasterID)
			conn.WriteJSON(WsPayload{Action: event.Action, Data: event.Data})

		case "message":
			// get clients conn from map using broadcasterID
			conn := m.Clients.Get(event.BroadcasterID)
			// Send message to client
			conn.WriteJSON(WsPayload{Action: event.Action, Data: event.Data})

		case "disconnect":
			// Handle client disconnect
			m.Clients.Remove(event.BroadcasterID)
		}
	}

}
