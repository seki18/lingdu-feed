# Lingdu Feed

A community feed platform with content discovery, social interactions, and a personalized recommendation algorithm. Built with **Go (Gin)** + **Next.js (React 19)** and **PostgreSQL**.

---

## Features

- **User System** — Registration, login (JWT), profile management, follow / unfollow
- **Content Feed** — Hybrid recommendation algorithm blending trending, recent, and social signals
- **Smart Interaction Tracking** — Feeds deduplicate seen content via fine-grained status states (delivered → displayed → clicked)
- **Posts & Comments** — Create, update, delete posts with nested reply-style comments
- **Engagement** — Likes (praises), collections (bookmarks), and view counts
- **Dual Auth** — Hard auth (required) + soft auth (optional, for guest browsing with degraded personalization)
- **History & Collections** — Personal reading history and saved posts

---

## Project Architecture

```
lingdu-feed/
├── backend/                    # Go API server (Gin + sqlx + PostgreSQL)
│   ├── cmd/main.go             # Entry point, route registration, CORS
│   ├── config/config.go        # Env-based configuration
│   ├── internal/
│   │   ├── common/             # Shared: DB pool, error codes, response envelope, pagination
│   │   ├── handler/            # HTTP handlers (thin, delegates to services)
│   │   ├── middleware/         # AuthMiddleware (JWT required) & SoftAuthMiddleware (optional)
│   │   ├── model/              # DB models & request/response DTOs
│   │   ├── repository/         # Data access layer (raw SQL via sqlx)
│   │   ├── router/             # Route grouping & middleware binding
│   │   ├── service/            # Business logic layer
│   │   └── utils/              # JWT helpers, Gin bind utilities
│   ├── migrations/             # SQL migration files
│   ├── go.mod
│   └── go.sum
├── frontend/                   # Next.js 16 + React 19 + TypeScript + Tailwind CSS 4
│   ├── src/
│   │   ├── app/                # App Router pages: feed, posts/[id], users/[id], etc.
│   │   ├── components/
│   │   │   ├── auth/           # LoginModal
│   │   │   ├── comment/        # CommentSection
│   │   │   ├── layout/         # Header, PostBlock, PostCard
│   │   │   └── ui/             # Loading, Toast (context-based notification)
│   │   ├── lib/                # apiFetch wrapper, auth token helpers
│   │   └── types/              # TypeScript interfaces for posts, comments, users
│   └── package.json
├── docker-compose.yml          # PostgreSQL 16 container
└── README.md
```

---

## Design Decisions

### Layered Backend (Handler → Service → Repository)

Each domain entity follows a strict three-layer separation:

| Layer | Responsibility | Rules |
|-------|---------------|-------|
| **Handler** | Parse request, call service, send response | No business logic |
| **Service** | Business logic, orchestration, feed algorithm | No direct SQL |
| **Repository** | Raw SQL queries via `sqlx` | No business logic |

This keeps the codebase testable and each layer replaceable independently.

### Feed Recommendation Algorithm

The `/feed/recommend` endpoint blends three signal sources:

1. **Recommend (≥50%)** — Posts picked by the interaction-status recommendation heuristics
2. **Recent (~33%)** — Fresh content from all users, with already-seen/dismissed posts filtered out
3. **Following (~17%)** — Posts from followed users

The algorithm ensures users always see a mix of discovery + familiarity, with deduplication via the `interaction_status` table.

### Soft Auth vs Hard Auth

Instead of requiring login for every endpoint, endpoints use one of two middlewares:

- **`AuthMiddleware`** — Rejects with `401` if no valid JWT (used for create/update/delete, personal feeds)
- **`SoftAuthMiddleware`** — Sets `user_id = -1` if unauthenticated, allowing guest access with best-effort personalization (used for public feeds, user profiles)

### Interaction Status Tracking

Every feed delivery records a status in `interaction_status`:
- `Unknown (0)` → `Delivered (1)` → `Displayed (2)` → `Clicked (3)`

This powers feed deduplication (never show the same post twice in a session) and enables future recommendation tuning based on engagement signals.

### Standardized API Envelope

All endpoints return:

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

Error codes follow a structured hierarchy (`400xx` client errors, `401xx` auth, `500xx` server errors). The original Go error is logged server-side but never leaked to the client.

### Counter Caching (Write-Through)

Like/praise, comment, collection, and follow counts are maintained as denormalized columns on the parent table. Every insert/delete on the child table triggers an atomic `UPDATE SET count = count ± 1` on the parent, so feeds never need expensive `COUNT(*)` joins.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend framework | [Gin](https://github.com/gin-gonic/gin) |
| Database driver | [sqlx](https://github.com/jmoiron/sqlx) + [lib/pq](https://github.com/lib/pq) |
| Auth | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) + bcrypt |
| Frontend framework | [Next.js 16](https://nextjs.org/) (App Router) |
| UI | [React 19](https://react.dev/) + [Tailwind CSS 4](https://tailwindcss.com/) |
| Language | TypeScript (strict) |
| Database | PostgreSQL 16 (via Docker Compose) |

---

## Getting Started

### Prerequisites

- Go 1.26+
- Node.js 20+
- Docker & Docker Compose

### 1. Start PostgreSQL

```bash
docker-compose up -d
```

### 2. Run database migrations

Apply the SQL files in `backend/migrations/` to your PostgreSQL instance.

### 3. Configure environment

Create a `.env` file in `backend/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=community
```

### 4. Start the backend

```bash
cd backend
go mod tidy
go run cmd/main.go
# Server runs on http://localhost:18080
```

### 5. Start the frontend

```bash
cd frontend
npm install
npm run dev
# App runs on http://localhost:3000
```

---

## API Overview

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/auth/register` | None | Create account |
| `POST` | `/auth/login` | None | Login, returns JWT |
| `GET` | `/users/me` | Required | Current user profile |
| `GET` | `/users/:id` | Soft | View any user profile |
| `PUT` | `/users` | Required | Update username |
| `GET` | `/feed/recommend` | Soft | Hybrid recommendation feed |
| `GET` | `/feed/following` | Required | Posts from followed users |
| `GET` | `/feed/author/:user_id` | None | Posts by a specific user |
| `GET` | `/feed/history` | Required | User's reading history |
| `GET` | `/feed/collections` | Required | User's saved posts |
| `POST` | `/post` | Required | Create a post |
| `PUT` | `/post` | Required | Update a post |
| `DELETE` | `/post` | Required | Delete a post |
| `POST` | `/comments` | Required | Add a comment |
| `DELETE` | `/comments` | Required | Delete a comment |
| `POST` | `/praises` | Required | Like a post |
| `DELETE` | `/praises` | Required | Unlike a post |
| `POST` | `/collections` | Required | Bookmark a post |
| `DELETE` | `/collections` | Required | Remove bookmark |
| `POST` | `/follows` | Required | Follow a user |
| `DELETE` | `/follows` | Required | Unfollow a user |
| `POST` | `/interaction-statuses` | Required | Record feed interaction |

---

## TODO

### Features

- [ ] **Image Upload** — Support post images (single/multiple) with file upload API and frontend preview
- [ ] **Search** — Full-text search across posts, users, and comments

### Optimizations

- [ ] **Caching Layer** — Introduce Redis for hot feed data and session caching to reduce DB load
- [ ] **Cloud Migration** — Migrate static assets to cloud object storage (S3/OSS) and deploy services to a cloud platform
