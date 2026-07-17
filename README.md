# Cinema Ticket Booking System

A real-time cinema ticket booking system built with Go (Gin), Vue 3, MongoDB, Redis, and RabbitMQ.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Frontend (Vue 3 + shadcn)                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐   │
│  │ TanStack │  │  Pinia   │  │  Router  │  │   SSE Client      │   │
│  │  Query   │  │  Stores  │  │          │  │ (EventSource)     │   │
│  └──────────┘  └──────────┘  └──────────┘  └────────┬──────────┘   │
│                                                      │              │
└──────────────────────────────────────────────────────┼──────────────┘
                           │ HTTP (REST)               │ SSE stream
                           ▼                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Backend (Go / Gin)                             │
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐   │
│  │ Handler  │→ │ Usecase  │→ │   Repo   │  │    SSE Handler    │   │
│  │  Layer   │  │  Layer   │  │  Layer   │  │ (Event Source)    │   │
│  └──────────┘  └──────────┘  └────┬─────┘  └────────┬──────────┘   │
│                                    │                  │             │
│  ┌─────────────────────────────────┼──────────────────┼───────────┐ │
│  │                                 ▼                  ▼           │ │
│  │                      ┌────────────────────┐                    │ │
│  │                      │   RabbitMQ (topic)  │                    │ │
│  │                      │  seat.* / booking.* │◄─────────────────── │ │
│  │                      └─────┬──────────────┘                    │ │
│  │                            │ seat.events                        │ │
│  │                      ┌─────▼──────────────┐                    │ │
│  │                      │   SSE Consumer     │──► SSE Clients     │ │
│  │                      └────────────────────┘                    │ │
│  └─────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐   │
│  │   MongoDB    │    │    Redis     │    │      RabbitMQ        │   │
│  │  (Persistence)│    │ (Distributed │    │   (Event Bus / MQ)  │   │
│  │              │    │    Locks)    │    │                      │   │
│  └──────────────┘    └──────────────┘    └──────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Tech Stack Overview

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Backend | Go 1.22 + Gin | REST API |
| Frontend | Vue 3 + TypeScript | SPA |
| UI | shadcn-vue + Tailwind CSS | Component library |
| State | Pinia + TanStack Query | Client state + server cache |
| Code Gen | Kubb → TanStack Query hooks | OpenAPI → typed hooks |
| Database | MongoDB 7 | Data persistence |
| Cache/Lock | Redis 7 | Distributed locks |
| Event Bus | RabbitMQ | Async events, audit, notifications |
| Real-time | SSE (Server-Sent Events) | Seat updates & notifications |
| Auth | Google OAuth 2.0 + JWT | Authentication |
| API Docs | Swagger (swaggo) | `/swagger/*` |
| Deployment | Docker Compose | `docker compose up --build` |

## Booking Flow

### Step-by-step

```
1. User views seat map → GET /api/v1/showtimes/:id/seats
   - Returns seat grid with AVAILABLE / LOCKED / BOOKED status

2. User clicks AVAILABLE seat → POST /api/v1/bookings/lock
   - Backend acquires Redis lock: SET lock:seat:{showtime}:{seat} {user_id} NX EX 300
   - Creates Booking document with status LOCKED, expires_at = now + 5min
   - Updates Seat status to LOCKED in MongoDB
   - Publishes seat.locked event to RabbitMQ

3. RabbitMQ seat.events consumer → SSE handler broadcasts to all connected clients
   - All users viewing this showtime see seat turn YELLOW (LOCKED) in real-time

4. User has 5 minutes to pay
   - Countdown timer shown on frontend

5. Options:
   ✅ User clicks "Confirm Payment" → POST /api/v1/bookings/pay
      - Releases Redis lock
      - Updates Booking status → BOOKED
      - Updates Seat status → BOOKED
      - Publishes seat.booked + notification events

   ❌ Timeout (5 min) → Background cleanup goroutine
      - Finds expired LOCKED bookings every 30s
      - Releases Redis lock
      - Updates Booking status → EXPIRED
      - Updates Seat status → AVAILABLE
      - Publishes seat.released + timeout events

   🔓 User clicks "Cancel" → POST /api/v1/bookings/cancel
      - Same as timeout but status → CANCELLED
```

## Redis Lock Strategy

### Design

```
Key:   lock:seat:{showtime_id}:{seat_id}
Value: {user_id}
TTL:   300 seconds (5 minutes)

Acquire: SET key value NX EX 300
Release: Lua script (verify ownership before delete)
Extend:  Lua script (verify ownership before EXPIRE)
```

### Why This Design

| Decision | Rationale |
|----------|-----------|
| **NX + EX atomic** | Prevents race between checking + setting lock |
| **User ID as value** | Ensures only the locker can unlock (ownership) |
| **Lua script for unlock** | Atomic check-then-delete prevents releasing others' locks |
| **5-minute TTL** | Matches booking lock window; auto-releases if client crashes |
| **Crash recovery** | Background goroutine scans MongoDB for expired LOCKED bookings (defense-in-depth) |

### Handling Failures

- **Client crash during lock**: Redis auto-expires in 5 minutes → cleanup goroutine releases MongoDB seats
- **Redis crash/restart**: Locks are lost, but cleanup goroutine is the safety net — scans for stale LOCKED bookings
- **Concurrent requests**: Multiple users locking same seat → only one gets SET NX success → others get 409 Conflict
- **Network partition**: Backend fails before creating Booking but after locking → Redis auto-expires the lock

## Message Queue (RabbitMQ) Usage

RabbitMQ is our event backbone. We use a **topic exchange** (`cinema`) for routing flexibility.

### Exchange & Queues

| Exchange | Routing Key | Queue | Consumer |
|----------|------------|-------|----------|
| `cinema` | `seat.*` | `seat.events` | SSE fan-out to connected clients |
| `cinema` | `booking.*` , `system.error` | `notification` | Mock notification logger |
| `cinema` | `audit.*` | (direct MongoDB) | N/A (written synchronously for audit) |

### Events Published

| Event | When | Payload |
|-------|------|---------|
| `seat.locked` | Lock acquired | `{booking_id, showtime_id, seat_ids, user_id, expires_at}` |
| `seat.booked` | Payment confirmed | `{booking_id, showtime_id, seat_ids, user_id}` |
| `seat.released` | Cancel/timeout | `{booking_id, showtime_id, seat_ids}` |
| `booking.success` | Payment confirmed | `{booking_id, user_id, seat_ids}` |
| `booking.timeout` | Lock expired | `{booking_id, user_id}` |
| `system.error` | Lock failure | `{seat_id, error}` |

### Why RabbitMQ over Kafka/Redis Pub-Sub

- We already have Redis for locks, but Redis Pub-Sub has no message persistence → messages lost if consumer disconnects
- RabbitMQ provides durable queues, acknowledgements, and dead-letter exchanges
- More production-grade for audit-critical events

## How to Run

### Prerequisites

- Docker & Docker Compose v2
- Google OAuth 2.0 credentials (see Setup section)

### Quick Start

```bash
# 1. Clone and set up environment
cp .env.example .env
# Edit .env with your Google OAuth credentials

# 2. Start everything
docker compose up --build
```

This starts:
- **MongoDB** on port 27017
- **Redis** on port 6379
- **RabbitMQ** on port 5672 (management UI: http://localhost:15672)
- **Backend** on port 8080
- **Seed** (runs once, seeds 3 movies × 2 screens × 3 showtimes + admin user)
- **Frontend** on port 5173

### Access

- **Frontend**: http://localhost:5173
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **RabbitMQ UI**: http://localhost:15672 (guest/guest)

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project or select existing
3. Enable **Google+ API** / **OAuth 2.0**
4. Create OAuth 2.0 Client ID (Web application)
5. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
6. Copy Client ID and Client Secret to `.env`

### Admin User

The seed script creates an admin user with email `admin@example.com`. Add this email to `ADMIN_EMAILS` in `.env` to grant admin role on first login.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Backend port |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection |
| `MONGO_DB` | `cinema_ticket` | Database name |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672` | RabbitMQ connection |
| `GOOGLE_CLIENT_ID` | - | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | - | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | `http://localhost:8080/...` | OAuth callback URL |
| `JWT_SECRET` | `dev-secret...` | JWT signing key |
| `JWT_ACCESS_EXPIRY` | `15m` | Access token TTL |
| `JWT_REFRESH_EXPIRY` | `168h` | Refresh token TTL |
| `FRONTEND_URL` | `http://localhost:5173` | CORS origin |
| `ADMIN_EMAILS` | `` | Comma-separated admin emails |

## Assumptions & Trade-offs

### Assumptions

1. **One screen per showtime**: Each Showtime maps to 1 screen with 40 seats (8 rows × 5 cols). The seed script creates this layout. A real system would have multiple screens per showtime.
2. **Single seat selection**: Current implementation supports multiple seats per booking, but the frontend UX is optimized for single-seat interaction.
3. **Mock payment**: The "Pay" endpoint simulates payment without real gateway integration. It just updates booking status.
4. **Google OAuth primary auth**: Only Google login is implemented. No email/password or other providers.

### Trade-offs

| Decision | Trade-off |
|----------|-----------|
| **SSE over WebSocket** | Simpler (native `EventSource`), but unidirectional (server→client only). For bidirectional real-time features (chat), WebSocket would be better. |
| **Direct audit log writes** | Audit logs are written synchronously to MongoDB in the usecase. For higher throughput, we'd move this to RabbitMQ consumers. |
| **Redis for locks only** | Not using Redis for caching or pub-sub (despite already having it). This keeps responsibilities clean: Redis = locks, RabbitMQ = events. |
| **No pagination** | API endpoints return all results. For production, add pagination to booking/audit-log endpoints. |
| **JWT in SSE URL** | SSE uses `EventSource` which doesn't support custom headers. The token is passed as query parameter. For production, implement token in the URL or use WebSocket with auth. |
| **Monorepo with /backend + /frontend** | Simple to deploy with Docker Compose. In larger teams, separate repos with CI per service would be better. |
| **No typing for API responses** | Frontend uses `as any[]` casts. With Kubb generating typed hooks from Swagger, this would be eliminated. |

## Project Structure

```
cinema_ticket/
├── backend/
│   ├── cmd/server/main.go           # Entry point, DI wiring
│   ├── internal/
│   │   ├── domain/                  # Entities & repository interfaces
│   │   ├── repository/mongo/        # MongoDB implementations
│   │   ├── repository/redis/        # Redis lock implementation
│   │   ├── usecase/                 # Business logic
│   │   ├── handler/                 # HTTP handlers
│   │   ├── middleware/              # JWT auth + role guard
│   │   ├── events/                  # RabbitMQ publisher/consumer
│   │   └── router/                  # Route registration
│   ├── pkg/config/                  # Env-based config
│   ├── scripts/seed.go             # Seed data script
│   ├── docs/                        # Swagger generated
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/              # UI + feature components
│   │   ├── views/                   # Page components (6 views)
│   │   ├── stores/                  # Pinia stores
│   │   ├── composables/             # useSSE, useAuth
│   │   ├── generated/               # Kubb → TanStack hooks
│   │   ├── router/                  # Vue Router
│   │   └── lib/                     # Utilities
│   ├── kubb.config.ts               # OpenAPI code gen
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```
