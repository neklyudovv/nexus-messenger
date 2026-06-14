package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r gin.IRouter) {
	g := r.Group("/users")
	g.GET("", h.getAll)
	g.GET("/me", h.getMe)
	g.PATCH("/me", h.updateMe)
}

func (h *Handler) getAll(c *gin.Context) {
	users, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	u, err := h.svc.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) updateMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req struct {
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.AvatarURL != "" && !strings.HasPrefix(req.AvatarURL, "https://") && !strings.HasPrefix(req.AvatarURL, "http://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar_url must be a valid http/https URL"})
		return
	}
	u, err := h.svc.Update(userID, req.Username, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}
