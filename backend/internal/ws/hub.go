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
	clients  map[uint]map[*Client]struct{} // userID → set of clients (multiple tabs)
	channels map[uint]map[uint]struct{}    // channelID → set of userIDs
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[uint]map[*Client]struct{}),
		channels: make(map[uint]map[uint]struct{}),
	}
}

// Register adds c to the hub. Returns true if this is the user's first connection.
func (h *Hub) Register(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.userID] == nil {
		h.clients[c.userID] = make(map[*Client]struct{})
	}
	first := len(h.clients[c.userID]) == 0
	h.clients[c.userID][c] = struct{}{}
	return first
}

// Unregister removes c from the hub. Returns true if the user has no more connections.
func (h *Hub) Unregister(c *Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[c.userID]
	if !ok {
		return false
	}
	delete(set, c)
	if len(set) > 0 {
		return false
	}
	delete(h.clients, c.userID)
	for _, members := range h.channels {
		delete(members, c.userID)
	}
	return true
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
	h.BroadcastRawToChannel(channelID, data)
}

func (h *Hub) BroadcastRawToChannel(channelID uint, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID := range h.channels[channelID] {
		for c := range h.clients[userID] {
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
	for _, set := range h.clients {
		for c := range set {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}

// BroadcastToUsers sends raw JSON only to the specified user IDs (all their connections).
func (h *Hub) BroadcastToUsers(userIDs []uint, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, userID := range userIDs {
		for c := range h.clients[userID] {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}
