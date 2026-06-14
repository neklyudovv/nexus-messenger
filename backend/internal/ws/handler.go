package ws

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"nexus-messenger/backend/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub        *Hub
	userSvc    presenceService
	msgSvc     messageService
	channelSvc channelService
	jwtSecret  string
	upgrader   websocket.Upgrader
}

func NewHandler(hub *Hub, userSvc presenceService, msgSvc messageService,
	channelSvc channelService, jwtSecret string, allowedOrigins []string) *Handler {

	checkOrigin := func(r *http.Request) bool { return true }
	if len(allowedOrigins) > 0 {
		checkOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, o := range allowedOrigins {
				if strings.TrimSpace(o) == origin {
					return true
				}
			}
			return false
		}
	}

	return &Handler{
		hub:        hub,
		userSvc:    userSvc,
		msgSvc:     msgSvc,
		channelSvc: channelSvc,
		jwtSecret:  jwtSecret,
		upgrader:   websocket.Upgrader{CheckOrigin: checkOrigin},
	}
}

func (h *Handler) Register(r gin.IRouter) {
	r.GET("/ws", h.serve)
}

// serve upgrades the connection and waits for the first message {"type":"auth","token":"<jwt>"}.
func (h *Handler) serve(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	var authMsg struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(raw, &authMsg); err != nil || authMsg.Type != "auth" || authMsg.Token == "" {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "auth required"))
		conn.Close()
		return
	}

	claims, err := auth.ParseToken(authMsg.Token, h.jwtSecret)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "invalid token"))
		conn.Close()
		return
	}

	conn.SetReadDeadline(time.Time{})

	client := NewClient(h.hub, conn, claims.UserID, h.userSvc, h.msgSvc, h.channelSvc)
	go client.Run()
}
