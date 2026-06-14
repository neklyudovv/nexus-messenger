package main

import (
	"log"
	"strings"
	"time"

	"nexus-messenger/backend/config"
	"nexus-messenger/backend/internal/auth"
	"nexus-messenger/backend/internal/channel"
	"nexus-messenger/backend/internal/db"
	"nexus-messenger/backend/internal/docs"
	"nexus-messenger/backend/internal/message"
	"nexus-messenger/backend/internal/user"
	"nexus-messenger/backend/internal/workspace"
	"nexus-messenger/backend/internal/ws"
	"nexus-messenger/backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	postgres, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}

	redis, err := db.NewRedis(cfg)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	authSvc      := auth.NewService(postgres, redis, cfg)
	userSvc      := user.NewService(postgres, redis)
	workspaceSvc := workspace.NewService(postgres)
	channelSvc   := channel.NewService(postgres)
	messageSvc   := message.NewService(postgres)
	hub          := ws.NewHub()

	allowedOrigins := splitOrigins(cfg.CORSOrigins)
	authLimiter    := middleware.RateLimit(10, time.Minute)

	authH      := auth.NewHandler(authSvc, cfg)
	userH      := user.NewHandler(userSvc)
	workspaceH := workspace.NewHandler(workspaceSvc)
	channelH   := channel.NewHandler(channelSvc, hub.BroadcastToUsers)
	messageH   := message.NewHandler(messageSvc, channelSvc, hub.BroadcastRawToChannel)
	wsH        := ws.NewHandler(hub, userSvc, messageSvc, channelSvc, cfg.JWTSecret, allowedOrigins)

	r := gin.Default()
	r.Use(corsMiddleware(allowedOrigins))

	r.GET("/api/openapi.yaml", func(c *gin.Context) {
		c.Data(200, "application/yaml; charset=utf-8", docs.OpenAPISpec)
	})
	r.GET("/api/docs", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", docs.ScalarHTML)
	})

	api := r.Group("/api")
	authH.Register(api, authLimiter)

	protected := api.Group("", middleware.Auth(cfg.JWTSecret))
	userH.Register(protected)
	workspaceH.Register(protected)
	channelH.Register(protected)
	messageH.Register(protected)

	wsH.Register(r)

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func splitOrigins(raw string) []string {
	if raw == "*" || raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func corsMiddleware(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqOrigin := c.Request.Header.Get("Origin")
		if reqOrigin == "" {
			c.Next()
			return
		}

		var acao string
		if len(allowed) == 0 {
			acao = reqOrigin // dev: reflect all origins
		} else {
			for _, o := range allowed {
				if o == reqOrigin {
					acao = reqOrigin
					break
				}
			}
		}

		if acao != "" {
			c.Header("Access-Control-Allow-Origin", acao)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
