package workspace

import (
	"net/http"

	"nexus-messenger/backend/internal/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r gin.IRouter) {
	g := r.Group("/workspaces")
	g.POST("", h.create)
	g.GET("/me", h.getMine)
	g.POST("/join", h.join)
	g.GET("/:id", h.getOne)
	g.GET("/:id/members", h.getMembers)
	g.POST("/:id/invite/regenerate", h.regenerateInvite)
	g.PATCH("/:id/members/:user_id", h.updateMemberRole)
}

func (h *Handler) create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req struct {
		Name        string `json:"name"        binding:"required,min=2,max=100"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ws, err := h.svc.Create(req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ws)
}

func (h *Handler) getMine(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	workspaces, err := h.svc.GetMine(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workspaces)
}

func (h *Handler) join(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req struct {
		InviteCode string `json:"invite_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ws, err := h.svc.Join(req.InviteCode, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func (h *Handler) getOne(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, ok := httputil.ParseID(c, "id")
	if !ok {
		return
	}
	if !h.svc.IsMember(id, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	ws, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func (h *Handler) getMembers(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, ok := httputil.ParseID(c, "id")
	if !ok {
		return
	}
	if !h.svc.IsMember(id, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	members, err := h.svc.GetMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *Handler) updateMemberRole(c *gin.Context) {
	requesterID := c.MustGet("user_id").(uint)
	wsID, ok := httputil.ParseID(c, "id")
	if !ok {
		return
	}
	targetID, ok := httputil.ParseID(c, "user_id")
	if !ok {
		return
	}
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateMemberRole(wsID, requesterID, targetID, MemberRole(req.Role)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

func (h *Handler) regenerateInvite(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, ok := httputil.ParseID(c, "id")
	if !ok {
		return
	}
	code, err := h.svc.RegenerateInvite(id, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invite_code": code})
}
