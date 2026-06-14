package message

import (
	"encoding/json"
	"net/http"
	"time"

	"nexus-messenger/backend/internal/channel"
	"nexus-messenger/backend/internal/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc        *Service
	channelSvc *channel.Service
	broadcast  func(channelID uint, data []byte)
}

func NewHandler(svc *Service, channelSvc *channel.Service, broadcast func(uint, []byte)) *Handler {
	return &Handler{svc: svc, channelSvc: channelSvc, broadcast: broadcast}
}

func (h *Handler) Register(r gin.IRouter) {
	r.GET("/workspaces/:id/channels/:cid/messages", h.getHistory)
	r.DELETE("/messages/:id", h.delete)
}

func (h *Handler) getHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	channelID, ok := httputil.ParseID(c, "cid")
	if !ok {
		return
	}

	if !h.channelSvc.IsMember(channelID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var before time.Time
	if b := c.Query("before"); b != "" {
		before, _ = time.Parse(time.RFC3339Nano, b)
	}

	msgs, err := h.svc.GetHistory(channelID, before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (h *Handler) delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	msgID, ok := httputil.ParseID(c, "id")
	if !ok {
		return
	}
	channelID, err := h.svc.Delete(msgID, userID)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}
	if h.broadcast != nil && channelID != 0 {
		if data, e := json.Marshal(map[string]any{
			"type":       "message_deleted",
			"channel_id": channelID,
			"payload":    map[string]any{"id": msgID},
		}); e == nil {
			h.broadcast(channelID, data)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
