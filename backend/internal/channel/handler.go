package channel

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nexus-messenger/backend/internal/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc       *Service
	broadcast func([]uint, []byte)
}

func NewHandler(svc *Service, broadcast func([]uint, []byte)) *Handler {
	return &Handler{svc: svc, broadcast: broadcast}
}

func (h *Handler) Register(r gin.IRouter) {
	r.GET("/channels/dms", h.listDMs)

	g := r.Group("/workspaces/:id")
	g.Use(h.workspaceCtx)

	g.GET("/channels", h.list)
	g.POST("/channels", h.create)
	g.GET("/channels/:cid", h.getOne)
	g.DELETE("/channels/:cid", h.delete)
	g.GET("/channels/:cid/members", h.members)
	g.POST("/channels/:cid/members", h.addMember)
	g.POST("/dm/:user_id", h.openDM)
}

func (h *Handler) workspaceCtx(c *gin.Context) {
	v, err := strconv.Atoi(c.Param("id"))
	if err != nil || v <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		c.Abort()
		return
	}
	c.Set("workspace_id", uint(v))
	c.Next()
}

func wid(c *gin.Context) uint { return c.MustGet("workspace_id").(uint) }
func uid(c *gin.Context) uint { return c.MustGet("user_id").(uint) }

func (h *Handler) list(c *gin.Context) {
	channels, err := h.svc.GetForWorkspace(wid(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, channels)
}

func (h *Handler) listDMs(c *gin.Context) {
	dms, err := h.svc.GetAllDMs(uid(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dms)
}

func (h *Handler) create(c *gin.Context) {
	var req struct {
		Name        string      `json:"name"        binding:"required,min=2,max=100"`
		Description string      `json:"description"`
		Type        ChannelType `json:"type"        binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ch, err := h.svc.Create(wid(c), req.Name, req.Description, req.Type, uid(c))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if h.broadcast != nil {
		if memberIDs, err := h.svc.GetWorkspaceMemberIDs(wid(c)); err == nil {
			if data, err := json.Marshal(map[string]any{"type": "channel_created", "payload": ch}); err == nil {
				h.broadcast(memberIDs, data)
			}
		}
	}
	c.JSON(http.StatusCreated, ch)
}

func (h *Handler) getOne(c *gin.Context) {
	cid, ok := httputil.ParseID(c, "cid")
	if !ok {
		return
	}
	ch, err := h.svc.GetByID(cid)
	if err != nil || ch.WorkspaceID != wid(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}
	if ch.Type != TypePublic && !h.svc.IsMember(cid, uid(c)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, ch)
}

func (h *Handler) delete(c *gin.Context) {
	cid, ok := httputil.ParseID(c, "cid")
	if !ok {
		return
	}
	if err := h.svc.Delete(cid, uid(c)); err != nil {
		if err.Error() == "channel not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) members(c *gin.Context) {
	cid, ok := httputil.ParseID(c, "cid")
	if !ok {
		return
	}
	if !h.svc.IsMember(cid, uid(c)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	members, err := h.svc.GetMembers(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *Handler) addMember(c *gin.Context) {
	cid, ok := httputil.ParseID(c, "cid")
	if !ok {
		return
	}
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddMember(cid, uid(c), req.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

func (h *Handler) openDM(c *gin.Context) {
	targetID, ok := httputil.ParseID(c, "user_id")
	if !ok {
		return
	}
	ch, err := h.svc.OpenDM(wid(c), uid(c), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ch)
}
