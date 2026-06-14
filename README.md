# Nexus Messenger

A self-hosted, Slack-style corporate messenger. One deployment = one company workspace.  
Real-time messaging with channels, direct messages, and online presence.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)

## Features

- **Workspaces** — create or join via invite code; admin/member roles
- **Channels** — public channels, real-time broadcast when created
- **Direct messages** — global, not tied to any workspace
- **Online presence** — who's online, typing indicators
- **Auth** — JWT access tokens (15 min) + HttpOnly refresh cookies (7 days), auto-rotating
- **API docs** — OpenAPI 3.0 spec + interactive Scalar UI at `/api/docs`

## Stack

| Layer     | Technology                              |
|-----------|-----------------------------------------|
| Backend   | Go 1.26, Gin, gorilla/websocket         |
| Auth      | JWT HS256, bcrypt, Redis token store    |
| Database  | PostgreSQL 16 (GORM, cursor pagination) |
| Cache     | Redis 7 (presence, typing, tokens)      |
| Frontend  | Vue 3, Pinia, Tailwind CSS v4, Vite     |
| Infra     | Docker Compose, nginx                   |

## Quick start

```bash
# 1. Copy and fill in secrets
cp .env.example .env

# 2. Start everything
docker-compose up --build
```

Open **http://localhost** -> register -> create a workspace -> share the invite code.

## Architecture

```
Browser ──► nginx :80
              ├── /         Vue 3 SPA (static files)
              ├── /api/*    Go backend :8080 (REST)
              └── /ws       Go backend :8080 (WebSocket)
```

### Backend layout

```
backend/
├── cmd/server/        # main: wires up all services and starts Gin
├── config/            # typed env config with validation
├── middleware/        # JWT auth, IP-based rate limiter
└── internal/
    ├── auth/          # register, login,  refresh,  logout
    ├── user/          # profile (username, avatar)
    ├── workspace/     # workspaces,  members,  invite codes
    ├── channel/       # channels,  DMs,  membership
    ├── message/       # history (cursor pagination),  soft-delete + WS broadcast
    ├── ws/            # WebSocket hub + per-client pump
    ├── db/            # postgres + redis clients
    ├── httputil/      # shared HTTP handler utilities (path param parsing)
    └── docs/          # embedded OpenAPI spec + Scalar HTML
```

### WebSocket protocol

Connect to `GET /ws`, then send:

```json
{ "type": "auth", "token": "<access_token>" }
```

Token goes in the **first message**, not in the URL (prevents leaking to logs).

**Client -> server** event types: `join_channel`,  `leave_channel`,  `send_message`,  `typing`,  `ping`

**Server -> client** event types: `new_message`,  `message_deleted`,  `user_online`,  `user_offline`,  `typing`,  `pong`,  `channel_created`

## Configuration

All settings are read from `.env` in the working directory.

| Variable           | Default     | Required | Notes                                 |
|--------------------|-------------|----------|---------------------------------------|
| `JWT_SECRET`       | —           | ✅       | Generate: `openssl rand -hex 32`      |
| `POSTGRES_PASSWORD`| —           | ✅       |                                       |
| `POSTGRES_HOST`    | `localhost` |          |                                       |
| `POSTGRES_SSL_MODE`| `disable`   |          | Set `require` in production           |
| `REDIS_HOST`       | `localhost` |          |                                       |
| `CORS_ORIGINS`     | `*`         |          | Set to your domain in production      |
| `SECURE_COOKIES`   | `false`     |          | Set `true` when serving over HTTPS    |
| `JWT_ACCESS_TTL`   | `15m`       |          |                                       |
| `JWT_REFRESH_TTL`  | `168h`      |          | 7 days                                |

## Security highlights

- Passwords hashed with **bcrypt** (cost 10)
- Refresh tokens stored as **HttpOnly; SameSite=Lax** cookies — inaccessible to JavaScript
- Refresh tokens are **server-side revocable** via Redis (logout invalidates immediately)
- Auth endpoints **rate-limited** to 10 req/min per IP
- PostgreSQL and Redis are **not exposed** outside the Docker network
- nginx sets **CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy** headers
- JWT signing method is validated on parse (algorithm confusion protection)

## Tests

```bash
cd backend
go test ./...
```

## API reference

Start the stack and visit **http://localhost/api/docs** for the interactive Scalar UI.  
Raw spec: `GET /api/openapi.yaml`
