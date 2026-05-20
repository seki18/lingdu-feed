# Lingdu Feed

A community content feed platform built with **Go (Gin)** + **Next.js 16 (React 19)** + **PostgreSQL + Redis**.

This system is designed as a **feed-based recommendation backend**, focusing on:

- hybrid ranking feed generation
- state-based deduplication pipeline
- user behavior tracking signals
- scalable backend architecture (Handler → Service → Repository)

---

## Features

| Area | Capabilities |
|---|---|
| **Auth** | Register, login, JWT (hard auth + optional soft auth) |
| **Posts** | CRUD with title + content, owner-only edit/delete |
| **Comments** | Nested replies, cascade delete |
| **Social** | Like/unlike, favorite/unfavorite, follow/unfollow |
| **Feeds** | Hybrid feed system (recommend + following + history) with state-based deduplication |
| **Tracking** | State pipeline (delivered → exposed → clicked) for feed deduplication and ranking signals |
| **API** | Unified JSON envelope `{ code, message, data }` with pagination |

---

## Tech Stack

| Layer | Stack |
|---|---|
| Backend | Go 1.26 + Gin v1.12 + sqlx v1.4 + PostgreSQL 16 |
| Frontend | Next.js 16 + React 19 + TypeScript 5 + Tailwind CSS 4 |
| Auth | golang-jwt v5 + bcrypt |
| Infra | Docker Compose (PostgreSQL 16, Redis 7) |

---

## Quick Start

```bash
# 1. Start PostgreSQL
docker compose up -d postgres

# 2. Run migration
docker compose exec -T postgres psql -U admin -d community < backend/migrations/001_init.sql

# 3. Configure .env (in backend/)
echo "DB_HOST=localhost
DB_PORT=15432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=community" > backend/.env

# 4. Start backend (:18080)
cd backend && go run cmd/main.go

# 5. Start frontend (:3000)
cd frontend && npm install && npm run dev
```

---

## Project Structure

```
lingdu-feed/
├── docker-compose.yml
├── backend/
│   ├── cmd/main.go
│   ├── config/config.go
│   ├── migrations/001_init.sql
│   └── internal/
│       ├── common/        # DB pool, errors, response helpers
│       ├── handler/       # HTTP handlers → feed, post, social, follow, state, user
│       ├── middleware/    # AuthMiddleware, SoftAuthMiddleware
│       ├── model/         # DB models + request/response DTOs
│       ├── repository/    # Raw SQL (sqlx), one file per table
│       ├── router/        # Route groups & middleware binding
│       ├── service/       # Business logic orchestration
│       └── utils/         # JWT, Gin helpers
├── frontend/
│   └── src/
│       ├── app/           # App Router (feed, posts/[id], users/[id], history, collections)
│       ├── components/    # auth, comment, layout, ui
│       ├── lib/           # api.ts (typed fetch), auth.ts
│       └── types/         # Post, Comment, User interfaces
```

---

## Architecture Overview

### Backend: Handler → Service → Repository

This architecture enforces strict separation of concerns:

- **Handler**: HTTP layer (request parsing & response formatting only)
- **Service**: business orchestration (feed composition, interaction logic, state updates)
- **Repository**: raw SQL layer (sqlx, no business logic)

This design allows:
- independent scaling of business logic
- clear separation between API and domain logic
- easier extension for caching and event tracking

### Design Highlights

- Hybrid feed generation instead of single ranking strategy
- State-based deduplication instead of naive caching
- Batch tracking to reduce network overhead
- Separation of state tracking and future event-based analytics
- Cursor-based pagination for scalable feed loading

### Database Schema

```
users ──1:N── posts ──1:N── comments (self-ref reply_id)
  │              │
  │              ├── N:M likes (user_id, post_id)
  │              ├── N:M favorites (user_id, post_id)
  │              └── N:M states (user_id, post_id, status 0-3)
  │
  └── N:M follows (follower_id, following_id)
```

---

## API Reference

All routes use `/api` prefix.

### Auth

| Method | Path | Auth |
|---|---|---|
| `POST` | `/api/auth/register` | Public |
| `POST` | `/api/auth/login` | Public |

### User

| Method | Path | Auth |
|---|---|---|
| `GET` | `/api/users/:id` | SoftAuth |
| `PUT` | `/api/users/me/profile` | Auth |
| `PUT` | `/api/users/me/password` | Auth |

### Feed

| Method | Path | Auth |
|---|---|---|
| `GET` | `/api/feed/recommend` | SoftAuth |
| `GET` | `/api/feed/following` | Auth |
| `GET` | `/api/feed/users/:id` | SoftAuth |
| `GET` | `/api/feed/history` | Auth |
| `GET` | `/api/feed/favorites` | Auth |

Params: `request_type` (`initial`/`subsequent`), `current_ids`, `page`, `page_size`

### Post

| Method | Path | Auth |
|---|---|---|
| `GET` | `/api/posts/:id` | SoftAuth |
| `POST` | `/api/posts` | Auth |
| `PUT` | `/api/posts/:id` | Auth |
| `DELETE` | `/api/posts/:id` | Auth |

### Social

| Method | Path | Auth |
|---|---|---|
| `POST` | `/api/posts/:id/like` | Auth |
| `DELETE` | `/api/posts/:id/like` | Auth |
| `POST` | `/api/posts/:id/favorite` | Auth |
| `DELETE` | `/api/posts/:id/favorite` | Auth |
| `GET` | `/api/posts/:id/comments` | SoftAuth |
| `POST` | `/api/posts/:id/comments` | Auth |
| `DELETE` | `/api/comments/:id` | Auth |

### Follow

| Method | Path | Auth |
|---|---|---|
| `POST` | `/api/users/:id/follow` | Auth |
| `DELETE` | `/api/users/:id/follow` | Auth |
| `GET` | `/api/users/:id/following` | Public |
| `GET` | `/api/users/:id/followers` | Public |

### State

| Method | Path | Auth |
|---|---|---|
| `POST` | `/api/state/batch` | Auth |

### Response Envelope

```json
{ "code": 200, "message": "success", "data": { ... } }
```

Paginated: `{ "items": [...], "total": 42, "page": 1, "page_size": 20 }`

### Error Codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 40001 | Invalid parameter |
| 40002 | Password incorrect |
| 40003 | User not found |
| 40004 | Post not found |
| 40100 | Unauthorized |
| 40901 | Email already registered |
| 50000 | Internal server error |

---

## Feed Algorithm (Hybrid Ranking System)

The feed system uses a hybrid ranking strategy combining:

- content freshness
- user engagement signals
- social graph signals

### Scoring

```
score = recency × 0.1 + views × 3 + likes × 5 + favorites × 4 + comments × 4
```

### Composition

| Source | Share | Filter |
|---|---|---|
| Recommend | ≥50% | Top-N by score, no state filter |
| Recent | ~33% | Latest posts, state-filtered |
| Following | ~17% | From followed users, state-filtered |

### Degradation

If normal requests return insufficient posts, the system auto-degrades: refetches from the recommend pool *without* the state filter, allowing previously-seen posts to fill remaining slots. The `excludeIDs` list is always preserved to avoid duplicates on the current page.

### State Pipeline (Feed Deduplication System)

This pipeline ensures feed consistency, deduplication, and interaction signal collection across sessions.

```
Delivered (1) → Exposed (2) → Clicked (3)
```

Reported in batch (500ms debounce). View count increments on first click only.

---

## TODO

### Features

- [ ] **Image Upload** — Support post images (single/multiple) with file upload API and frontend preview
- [ ] **Search** — Full-text search across posts, users, and comments
- [ ] **User Audit Log** — Record user login/logout, page dwell time, and key actions for analytics
- [ ] **Observability** — Structured logging, request tracing, error alerting, and incident investigation toolchain

### Optimizations

- [ ] **Caching Layer** — Introduce Redis for hot feed data and session caching to reduce DB load
- [ ] **Cloud Migration** — Migrate static assets to cloud object storage (S3/OSS) and deploy services to a cloud platform

---

## Development

```bash
# Backend
cd backend && go run cmd/main.go     # :18080

# Frontend
cd frontend && npm run dev           # :3000

# Database shell
docker compose exec postgres psql -U admin -d community
```

