package ws

import (
	"encoding/json"
	"log"
	"time"

	"nexus-messenger/backend/internal/channel"
	"nexus-messenger/backend/internal/message"
	"nexus-messenger/backend/internal/user"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 25 * time.Second
	maxMsgSize = 4096
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	userID     uint
	userSvc    *user.Service
	msgSvc     *message.Service
	channelSvc *channel.Service
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uint,
	userSvc *user.Service, msgSvc *message.Service, channelSvc *channel.Service) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		userID:     userID,
		userSvc:    userSvc,
		msgSvc:     msgSvc,
		channelSvc: channelSvc,
	}
}

type incomingEvent struct {
	Type      string `json:"type"`
	ChannelID uint   `json:"channel_id"`
	Content   string `json:"content"`
}

func (c *Client) Run() {
	c.hub.Register(c)
	c.userSvc.SetOnline(c.userID)
	c.hub.BroadcastToAll(Event{Type: "user_online", UserID: c.userID})

	go c.writePump()
	c.readPump()

	c.hub.Unregister(c)
	c.userSvc.SetOffline(c.userID)
	c.hub.BroadcastToAll(Event{Type: "user_offline", UserID: c.userID})
	c.conn.Close()
}

func (c *Client) readPump() {
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.userSvc.SetOnline(c.userID)
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var ev incomingEvent
		if err := json.Unmarshal(raw, &ev); err != nil {
			continue
		}
		c.handle(ev)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case data, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, data)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handle(ev incomingEvent) {
	switch ev.Type {
	case "join_channel":
		if ev.ChannelID == 0 || !c.isMember(ev.ChannelID) {
			return
		}
		c.hub.JoinChannel(ev.ChannelID, c.userID)

	case "leave_channel":
		c.hub.LeaveChannel(ev.ChannelID, c.userID)

	case "send_message":
		if ev.Content == "" || ev.ChannelID == 0 || !c.isMember(ev.ChannelID) {
			return
		}
		msg, err := c.msgSvc.Create(ev.ChannelID, c.userID, ev.Content)
		if err != nil {
			log.Printf("create message: %v", err)
			return
		}
		c.hub.BroadcastToChannel(ev.ChannelID, Event{
			Type:      "new_message",
			ChannelID: ev.ChannelID,
			Payload:   msg,
		})

	case "typing":
		if ev.ChannelID == 0 || !c.isMember(ev.ChannelID) {
			return
		}
		c.userSvc.SetTyping(ev.ChannelID, c.userID)
		c.hub.BroadcastToChannel(ev.ChannelID, Event{
			Type:      "typing",
			ChannelID: ev.ChannelID,
			UserID:    c.userID,
		})

	case "ping":
		c.userSvc.SetOnline(c.userID)
		data, _ := json.Marshal(Event{Type: "pong"})
		select {
		case c.send <- data:
		default:
		}
	}
}

func (c *Client) isMember(channelID uint) bool {
	if !c.channelSvc.IsMember(channelID, c.userID) {
		c.sendError("access denied")
		return false
	}
	return true
}

func (c *Client) sendError(msg string) {
	data, _ := json.Marshal(Event{Type: "error", Payload: msg})
	select {
	case c.send <- data:
	default:
	}
}
