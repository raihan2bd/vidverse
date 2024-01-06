package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/helpers"
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

func (c *Clients) RemoveByConn(conn *websocket.Conn) {
	c.Lock()
	defer c.Unlock()
	for k, v := range c.m {
		if v == conn {
			delete(c.m, k)
		}
	}
}

func (c *Clients) Count() int64 {
	c.RLock()
	defer c.RUnlock()
	return int64(len(c.m))
}

func (c *Clients) Get(userID uint) *websocket.Conn {
	c.RLock()
	defer c.RUnlock()
	return c.m[userID]
}

type WsPayload struct {
	Action string          `json:"action"`
	Data   interface{}     `json:"data,omitempty"`
	Conn   *websocket.Conn `json:"-"`
}

type ErrorRes struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}

// handle websocket request
func (m *Repo) WSHandler(c *gin.Context) {

	// Upgrade HTTP connection to WebSocket
	conn, err := upgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Extract user ID from context or request
	tokenString := c.Query("token")
	if tokenString == "" {
		conn.WriteJSON(WsPayload{Action: "unauthorized", Data: ErrorRes{Error: "Unauthorized", Status: http.StatusUnauthorized}})
		conn.Close()
		return
	}

	// validate token
	token, err := helpers.DecodeToken(tokenString)
	if err != nil {
		conn.WriteJSON(WsPayload{Action: "unauthorized", Data: ErrorRes{Error: "Unauthorized", Status: http.StatusUnauthorized}})
		conn.Close()
		return
	}

	if !helpers.ValidateToken(token) {
		conn.WriteJSON(WsPayload{Action: "unauthorized", Data: ErrorRes{Error: "Unauthorized", Status: http.StatusUnauthorized}})
		conn.Close()
		return
	}

	userID := uint(token["sub"].(float64))

	// Add client to map
	m.Clients.Add(userID, conn)

	// Send connection event message
	m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "notifications", Data: []string{}}
	// conn.WriteJSON(WsPayload{Action: "connect", Data: "connected"})

	go func() {

		defer func() {
			if r := recover(); r != nil {
				log.Println("Error", fmt.Sprintf("%v", r))
			}
		}()

		var payload WsPayload

		for {
			err := conn.ReadJSON(&payload)
			if err != nil {
				// Handle error and send disconnect event
				m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Action: "disconnect"}
				continue
			}

			switch payload.Action {
			case "close":
				m.Clients.RemoveByConn(conn)
				conn.Close()
				fmt.Println("Connection closed")

			default:
				continue
			}

			payload.Conn = conn

			// Send message event
			m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: userID, Data: payload, Action: payload.Action, Conn: payload.Conn}
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
			if conn != nil {
				fmt.Println("sending notification", event.Data)
				err := conn.WriteJSON(WsPayload{Action: "a_new_notification", Data: event.Data})
				if err != nil {
					continue
				}
			} else {
				continue
			}

		case "notifications":
			// get clients conn from map using broadcasterID
			conn := m.Clients.Get(event.BroadcasterID)
			if conn == nil {
				continue
			} else {

				// Send message to client
				count, err := m.App.DBMethods.GetUnreadNotificationsCountByUserID(event.BroadcasterID)
				if err != nil {
					count = 0
				}
				err = conn.WriteJSON(WsPayload{Action: event.Action, Data: count})
				if err != nil {
					continue
				}
			}

		case "disconnect":
			// Handle client disconnect
			m.Clients.Remove(event.BroadcasterID)
		}
	}

}
