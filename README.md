# Lingdu Feed

A community content feed platform built with **Go (Gin)** + **Next.js 16 (React 19)** + **PostgreSQL + Redis**.

This system is designed as a **feed-based recommendation backend**, focusing on:

- hybrid ranking feed generation with scheduled score computation
- Redis caching for hot feed data and user follow lists
- cursor-based pagination across multiple recall sources
- state-based deduplication pipeline
- user behavior tracking signals (view, click-through-rate)
- scalable backend architecture (Handler → Service → Repository)

---

## Features

| Area | Capabilities |
|---|---|
| **Auth** | Register, login, JWT (hard auth + optional soft auth) |
| **Posts** | CRUD with title + content, owner-only edit/delete |
| **Comments** | Nested replies, cascade delete |
| **Social** | Like/unlike, favorite/unfavorite, follow/unfollow |
| **Cache** | Redis: ranking ZSET, candidate ZSET, follow SET, feeditem HASH, content STRING, stats HASH, consumed state SET; cache-first, write-back, and write-through patterns |
| **Feeds** | Hybrid feed system (recommend + recent + following) with cursor-based pagination and state-based deduplication |
| **Images** | Multi-image upload (up to 9 per post), S3 storage, automatic compression (max 1920px, JPEG Q=85), feed thumbnails |
| **Tracking** | State pipeline (delivered → exposed → clicked) for feed deduplication and ranking signals |
| **API** | Unified JSON envelope `{ code, message, data }` with pagination |

---

## Tech Stack

| Layer | Stack |
|---|---|
| Backend | Go 1.26 + Gin v1.12 + sqlx v1.4 + PostgreSQL 16 |
| Frontend | Next.js 16 + React 19 + TypeScript 5 + Tailwind CSS 4 |
| Cache | Redis 7 + go-redis/v9 |
| Storage | AWS S3 + aws-sdk-go-v2 |
| Auth | golang-jwt v5 + bcrypt |
| Infra | Docker Compose (PostgreSQL 16, Redis 7) |

---

## Quick Start

```bash
# 1. Start PostgreSQL & Redis
docker compose up -d postgres redis

# 2. Run migration
docker compose exec -T postgres psql -U admin -d community < backend/migrations/001_init.sql

# 3. Create config.yaml from template (in backend/)
cp backend/config.example.yaml backend/config.yaml
# Then edit backend/config.yaml with your actual credentials

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
│       ├── common/        # DB pool, Redis client, S3 client, errors, response helpers
│       ├── handler/       # HTTP handlers → feed, post, social, follow, state, user, upload
│       ├── middleware/    # AuthMiddleware, SoftAuthMiddleware
│       ├── model/         # DB models + request/response DTOs (post, image, user, …)
│       ├── cache/         # Redis business logic (ranking, candidate, follow, stats, feeditem, content, state)
│       ├── repository/    # Raw SQL (sqlx), one file per table (post, image, follow, …)
│       ├── router/        # Route groups & middleware binding
│       ├── scheduler/     # Background tasks (score recalculation)
│       ├── service/       # Business logic orchestration
│       ├── storage/       # Image compression + S3 upload
│       └── utils/         # JWT, filter helpers, Gin helpers
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

- Scheduled score computation with CTR metric, decoupled from query time
- Redis hybrid state store (cache + write-back stats + persistent state)
- Hybrid feed generation from three recall sources (recommend / recent / following)
- Cursor-based pagination across all recall sources for scalable feed loading
- State-based deduplication instead of naive caching
- Batch tracking to reduce network overhead
- Separation of state tracking and future event-based analytics

### Database Schema

```
users ──1:N── posts ──1:N── comments (self-ref reply_id)
  │              │
  │              ├── 1:N── post_images (image_url, sort_order)
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

Params: `request_type` (`initial`/`subsequent`), `cursor` (pagination id), `page`, `page_size`

### Post

| Method | Path | Auth |
|---|---|---|
| `GET` | `/api/posts/:id` | SoftAuth |
| `POST` | `/api/posts` | Auth |
| `PUT` | `/api/posts/:id` | Auth |
| `DELETE` | `/api/posts/:id` | Auth |
| `POST` | `/api/posts/:id/images` | Auth |

### Upload

| Method | Path | Auth |
|---|---|---|
| `POST` | `/api/upload` | Auth |

Body: multipart/form-data with `file` field and optional `post_id` field.
Returns: `{ "url": "https://..." }`

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

The feed system uses a hybrid ranking strategy combining content freshness,
user engagement signals, and social graph signals.

### Score Formula (computed periodically by scheduler)

A normalized score ∈ [0, 1] is recalculated every minute:

| Component | Weight | Formula | Purpose |
|---|---|---|---|
| Recency | 15% | `EXP(-age / 7d-half-life)` | Time decay, 7-day half-life |
| Popularity | 35% | `tanh(views / 200)` | Absolute view volume |
| CTR | 20% | `views / expose_count` | Click-through rate |
| Likes | 15% | `tanh(likes / 50)` | Like engagement |
| Comments | 10% | `tanh(comments / 30)` | Discussion engagement |
| Favorites | 5% | `tanh(favorites / 30)` | Save/bookmark rate |

The scoring is decoupled from query time: on startup, a full-table update runs;
subsequently, posts whose stats changed within 24 hours or whose age is within
21 days (~3 recency half-lives) are recalculated each tick. This ensures
the recency decay term naturally degrades older content without requiring
new interactions.

### Redis Architecture

The system uses Redis as a hybrid state store — combining pure caching (cache-first with DB fallback),
write-back buffering (stats HINCRBY → batch sync), and persistence for user state (consumed SET with sliding TTL).

| Key | Structure | Contents | TTL / Cap | Write Strategy |
|---|---|---|---|---|
| `ranking` | ZSET | Top-N posts by score | refreshed on scheduler tick | scheduler → DB → Redis |
| `candidate` | ZSET | Posts from last 3 days by `created_time` | pruned: entries older than 3 days | post creation → ZADD + ZREMRANGEBYSCORE |
| `follow:<uid>` | SET | User's following user IDs | 24 hours, invalidated on follow/unfollow | follow/unfollow → SADD/SREM or DEL |
| `stats:{id}` | HASH | like/comment/favorite/view/expose count + score | 1 hour | HINCRBY on mutation, 30s batch sync to DB |
| `feeditem:{id}` | HASH | id, user_id, username, title, created_time | 1 hour | post create/update → HSET, post delete → DEL |
| `content:{id}` | STRING | post content text | 1 hour | post create/update → SETEX, post delete → DEL |
| `consumed:{uid}` | SET | post IDs the user has seen (>delivered) | 30 min (sliding) | BatchUpsertState → SADD + EXPIRE |

### Read / Write Patterns

- **Cache-first** (`ranking`, `candidate`, `follow`, `feeditem`, `content`): Try Redis first; on miss, fall back to DB and backfill cache.
- **Write-back** (`stats`): Mutations hit Redis HINCRBY immediately. A background goroutine syncs all dirty stats to the `post_stats` table every 30 seconds.
- **Write-through** (`consumed`): DB write (state) and Redis SADD happen simultaneously. Redis serves as a read-optimized replica. Sliding TTL keeps the SET alive for active users; inactive users' SET expires after 30 minutes and is rebuilt from DB on next request.

### Design Rationale

- `FeedItem` and `content`: reduce feed/detail queries from hitting the database. Feed item metadata is small (HASH) and rarely changes; content is large (TEXT) and only needed on the detail page.
- `stats`: absorbs high-frequency like/comment/view counter mutations using Redis atomic HINCRBY, eliminating per-request DB UPDATEs.
- `consumed`: replaces per-request DB queries in `FilterPostIDs` with an in-memory SET difference, significantly reducing feed assembly latency. Redis acts as a hot replica of persistent user consumption state.
- `ranking`: serves as a hot index; Top-N posts cover many pages of cursor-based pagination beyond the first page.
- `candidate`: uses time-based pruning (3 days) instead of count-based capping (20), ensuring cursor pagination works correctly across all pages within the window.

### Feed Composition

Three recall sources are combined per request:

| Source | Share | Sort | Cursor |
|---|---|---|---|
| Recommend | ≥50% | `score DESC, id DESC` | `(score, id) cursor` |
| Recent | ~33% | `created_time DESC` | `id < cursorID` |
| Following | ~17% | `created_time DESC` | `id < cursorID` |

All sources share a single `id` cursor for reliable, lossless pagination
across dynamic data.

### Degradation

If normal requests return insufficient posts, the system auto-degrades:
refetches from the recommend pool *without* the state filter, allowing
previously-seen posts to fill remaining slots.

---

## TODO

### Phase 1 — Infra & Delivery

- [x] **Image Upload (S3)** — Single/multiple image upload (max 9 per post), automatic compression, feed thumbnails
- [ ] **Deployment** — Containerize services, deploy to cloud platform with auto-scaling
- [ ] **CI/CD** — Automated build, test, lint, and deploy pipeline
- [ ] **Observability** — Structured logging (ELK), distributed tracing (Jaeger/Otel), metrics (Prometheus + Grafana), alerting

### Phase 2 — Features & Engagement

- [ ] **Search (Elasticsearch)** — Full-text search across posts, users, and comments with relevance ranking
- [ ] **Notification** — Push/in-app notifications for likes, comments, follows, and system events
- [ ] **Timeline Merge Deepening** — Smarter hybrid feed merging with decay curves, diversity constraints, and cold-start boosting

### Phase 3 — Scale

- [ ] **Kafka** — Event-driven architecture for async processing (feed delivery, notifications, analytics)

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

