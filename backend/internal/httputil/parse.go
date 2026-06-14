package httputil

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseID extracts a positive integer path parameter and writes a 400 response on failure.
func ParseID(c *gin.Context, param string) (uint, bool) {
	v, err := strconv.Atoi(c.Param(param))
	if err != nil || v <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + param})
		return 0, false
	}
	return uint(v), true
}
