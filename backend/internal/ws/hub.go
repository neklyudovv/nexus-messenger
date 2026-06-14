package ws

import (
	"encoding/json"
	"sync"
)

type Event struct {
	Type      string `json:"type"`
	ChannelID uint   `json:"channel_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	Payload   any    `json:"payload,omitempty"`
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[uint]*Client          // userID → client
	channels map[uint]map[uint]struct{} // channelID → set of userIDs
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[uint]*Client),
		channels: make(map[uint]map[uint]struct{}),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.userID] = c
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c.userID)
	for _, members := range h.channels {
		delete(members, c.userID)
	}
}

func (h *Hub) JoinChannel(channelID, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.channels[channelID] == nil {
		h.channels[channelID] = make(map[uint]struct{})
	}
	h.channels[channelID][userID] = struct{}{}
}

func (h *Hub) LeaveChannel(channelID, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if members, ok := h.channels[channelID]; ok {
		delete(members, userID)
	}
}

func (h *Hub) BroadcastToChannel(channelID uint, event Event) {
	data, _ := json.Marshal(event)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID := range h.channels[channelID] {
		if c, ok := h.clients[userID]; ok {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastToAll(event Event) {
	data, _ := json.Marshal(event)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		select {
		case c.send <- data:
		default:
		}
	}
}

// BroadcastRawToChannel sends pre-marshalled JSON to every client subscribed to a channel.
func (h *Hub) BroadcastRawToChannel(channelID uint, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID := range h.channels[channelID] {
		if c, ok := h.clients[userID]; ok {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}

// BroadcastToUsers sends raw JSON only to the specified user IDs.
func (h *Hub) BroadcastToUsers(userIDs []uint, data []byte) {
	set := make(map[uint]struct{}, len(userIDs))
	for _, id := range userIDs {
		set[id] = struct{}{}
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID, c := range h.clients {
		if _, ok := set[userID]; ok {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}
