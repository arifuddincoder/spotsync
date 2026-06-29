# 🚗 SpotSync – Smart Parking & EV Charging Reservation API

> A centralized parking management platform for busy airports and shopping malls.
> Alongside general parking zones, it safely handles the reservation of high-demand **EV charging spots** without any race conditions.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Echo](https://img.shields.io/badge/Echo-v5-00B5AD)](https://echo.labstack.com/)
[![GORM](https://img.shields.io/badge/GORM-PostgreSQL-336791)](https://gorm.io/)
[![License](https://img.shields.io/badge/license-MIT-green)]()

---

## 🔗 Live Links

| Resource | Link |
| --- | --- |
| 🐙 **GitHub Repo** | https://github.com/arifuddincoder/spotsync |
| 🚀 **Live API** | https://spotsync-b9vo.onrender.com |
| 🎤 **Interview Video** | https://drive.google.com/file/d/17KuOrfNjdlCQju92CTl4RwkeVS1xmHWw/view?usp=sharing |

> Base URL: `https://spotsync-b9vo.onrender.com/api/v1`

---

## 📌 Project Overview

SpotSync is a **REST API backend** that lets drivers browse parking/EV zones, reserve a spot, and cancel their own reservations. Admins can create/update/delete zones, set pricing, and view every reservation in the system.

The core challenge of the project is the **"EV Spot Bottleneck"** — making sure a zone never goes over capacity even when many users try to reserve the same spot at the exact same moment. This is solved using a **GORM Transaction + Row-Level Lock (`FOR UPDATE`)** (details below).

---

## ✨ Features

- 🔐 **JWT Authentication** — register/login, token-based sessions, bcrypt password hashing.
- 👥 **Role-Based Access Control** — `driver` and `admin` roles, enforced via middleware.
- 🅿️ **Parking Zone Management** — full CRUD (admin only) plus public listing.
- 📊 **Dynamic Availability** — each zone's `available_spots` is calculated live (total capacity − active reservation count).
- 🚦 **Concurrency-Safe Reservation** — row locking prevents the last-EV-spot race condition.
- 🚗 **Duplicate Plate Guard** — the same license plate cannot hold two active reservations at once (transaction check + unique partial index).
- 🧱 **Clean / Layered Architecture** — handler → service → repository are fully separated.
- 📦 **Standardized JSON Responses** — consistent success/error format across the whole API.
- 🌱 **Admin Auto-Seeding** — a default admin is created from `.env` on server boot.

---

## 🛠️ Tech Stack

| Technology | Usage |
| --- | --- |
| **Go (Golang)** | Core programming language |
| **Echo v5** (`github.com/labstack/echo/v5`) | HTTP web framework, routing, middleware |
| **GORM** (`gorm.io/gorm`) | ORM with PostgreSQL driver |
| **PostgreSQL** | Relational database (NeonDB/Supabase) |
| **validator/v10** | DTO struct validation |
| **golang-jwt/jwt/v5** | JWT token signing & verification |
| **bcrypt** (`golang.org/x/crypto`) | Password hashing |
| **godotenv** | Loading `.env` variables |
| **Air** | Hot-reload during local development |

---

## 🏛️ Architecture (Domain-Based Clean Architecture)

This project uses a **feature/domain-based** layout. The application is split into domains (`user`, `zone`, `reservation`), and **each domain contains its own layers** (dto, entity, handler, service, repository). This keeps all the code for a single feature in one place and makes it easy to maintain.

### 📂 Folder Structure

```
spotsync/
├── cmd/
│   └── main.go                  # Entry point — DI wiring, migration, seeding
│
├── internal/
│   ├── auth/
│   │   └── jwt.go               # JWT service (generate/validate token)
│   │
│   ├── config/
│   │   ├── config.go            # Loads .env into a Config struct
│   │   └── db.go                # GORM + PostgreSQL connect, connection pooling
│   │
│   ├── domain/                  # 🧩 All business domains live here
│   │   ├── user/
│   │   │   ├── dto/
│   │   │   │   ├── request.go   # RegisterRequest, LoginRequest
│   │   │   │   └── response.go  # UserResponse, LoginResponse
│   │   │   ├── entity.go        # User GORM model + password hash/check
│   │   │   ├── handler.go       # HTTP layer (bind, validate, response)
│   │   │   ├── service.go       # Business logic (hash, JWT, rules)
│   │   │   ├── repository.go    # DB access (CRUD)
│   │   │   ├── seed.go          # Default admin seeding
│   │   │   └── register.go      # Route registration
│   │   │
│   │   ├── zone/
│   │   │   ├── dto/ (request.go, response.go)
│   │   │   ├── entity.go        # ParkingZone model
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go    # includes available_spots subquery
│   │   │   └── register.go
│   │   │
│   │   └── reservation/
│   │       ├── dto/ (request.go, response.go)
│   │       ├── entity.go        # Reservation model (User, Zone relations)
│   │       ├── handler.go
│   │       ├── service.go       # ownership/cancel rules
│   │       ├── repository.go    # 🔒 Transaction + Row Lock lives here
│   │       ├── migrate.go       # Unique partial index migration
│   │       └── register.go
│   │
│   ├── httpresponse/
│   │   └── response.go          # Standard Success/Error response helpers
│   │
│   ├── middlewares/
│   │   └── auth.go              # AuthMiddleware + RequireRole
│   │
│   └── server/
│       └── http.go              # Echo init, CORS, validator, route mount
│
├── .air.toml                    # Air hot-reload config
├── .env.example                 # Sample environment variables
├── .gitignore
├── go.mod
└── go.sum
```

### 🔄 How the Layers Interact (Request Flow)

Every request follows the path below — no layer steps outside its own responsibility:

```
   HTTP Request
        │
        ▼
┌─────────────────┐   Router (register.go) routes to the correct handler
│   Middleware    │   AuthMiddleware → verify JWT, inject user into context
│                 │   RequireRole    → check role (admin/driver)
└────────┬────────┘
         ▼
┌─────────────────┐   • Bind request body into a DTO
│    Handler      │   • Validate with the validator
│   (HTTP layer)  │   • Read user_id/role from context
│                 │   • Call the Service, return a JSON response
└────────┬────────┘   ❌ No DB queries here
         ▼
┌─────────────────┐   • Enforce business rules (capacity, ownership)
│    Service      │   • Hash passwords, generate JWTs
│ (Business Logic)│   • Map Entity ↔ DTO
└────────┬────────┘   ❌ No HTTP/echo logic here
         ▼
┌─────────────────┐   • All GORM operations (CRUD)
│   Repository    │   • Transactions, Row Locks, subqueries
│  (Data Access)  │
└────────┬────────┘
         ▼
   PostgreSQL (via GORM)
```

**Core principles:**
- The **Handler** never talks to the database directly.
- The **Repository** knows nothing about HTTP/echo.
- The **Service** is the brain in the middle — all rules live here.
- Using **DTOs** ensures GORM models are never exposed directly to the API (so passwords never leak).

### 🧷 Dependency Injection (`main.go`)

Layers are wired manually — Repository → Service → Handler:

```go
repo := NewRepository(db)
svc  := NewService(repo, jwtService)
h    := NewHandler(svc)
```

On boot, `cmd/main.go` does: load config → connect DB → AutoMigrate → reservation index migration → seed admin → start server.

---

## 🗄️ Database Schema

### `users`
| Field | Type | Note |
| --- | --- | --- |
| id | uint (PK) | auto increment |
| name | varchar(100) | required |
| email | varchar(255) | unique, required |
| password | varchar(100) | bcrypt hash (never returned in responses) |
| role | varchar(20) | `driver` / `admin`, default `driver` |
| created_at / updated_at / deleted_at | timestamp | GORM auto (soft delete) |

### `parking_zones`
| Field | Type | Note |
| --- | --- | --- |
| id | uint (PK) | auto increment |
| name | varchar(100) | required |
| type | varchar(20) | `general` / `ev_charging` / `covered` |
| total_capacity | int | > 0 |
| price_per_hour | decimal(10,2) | > 0 |
| timestamps | timestamp | GORM auto |

### `reservations`
| Field | Type | Note |
| --- | --- | --- |
| id | uint (PK) | auto increment |
| user_id | uint (FK → users) | indexed |
| zone_id | uint (FK → parking_zones) | indexed |
| license_plate | varchar(15) | required |
| status | varchar(20) | `active` / `completed` / `cancelled`, default `active` |
| timestamps | timestamp | GORM auto |

> 🔒 **Extra constraint:** A **unique partial index** on `reservations` ensures a given
> `license_plate` can hold only one `active` reservation at a time.
> ```sql
> CREATE UNIQUE INDEX uniq_active_license_plate
> ON reservations (license_plate)
> WHERE status = 'active' AND deleted_at IS NULL;
> ```

---

## 🌐 API Endpoints

Every response follows this standard format:
```json
{ "success": true,  "message": "...", "data": { } }
{ "success": false, "message": "...", "errors": "..." }
```

### 🔹 Auth Module
| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| POST | `/api/v1/auth/register` | Public | Register a new user |
| POST | `/api/v1/auth/login` | Public | Login + JWT token |

### 🔹 Zone Module
| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| GET | `/api/v1/zones` | Public | All zones + `available_spots` |
| GET | `/api/v1/zones/:id` | Public | Single zone detail |
| POST | `/api/v1/zones` | Admin | Create a new zone |
| PUT | `/api/v1/zones/:id` | Admin | Update a zone (partial) |
| DELETE | `/api/v1/zones/:id` | Admin | Delete a zone |

### 🔹 Reservation Module
| Method | Endpoint | Access | Description |
| --- | --- | --- | --- |
| POST | `/api/v1/reservations` | Auth | Reserve a spot (⚠️ concurrency-safe) |
| GET | `/api/v1/reservations/my-reservations` | Auth | Current user's reservations |
| DELETE | `/api/v1/reservations/:id` | Auth | Cancel own reservation |
| GET | `/api/v1/reservations` | Admin | All reservations in the system |

<details>
<summary><b>📥 Sample Request/Response (click to expand)</b></summary>

**Register** — `POST /api/v1/auth/register`
```json
{ "name": "John Doe", "email": "john@spotsync.com", "password": "securePass123", "role": "driver" }
```

**Login Response** — `200 OK`
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiI...",
    "user": { "id": 1, "name": "John Doe", "email": "john@spotsync.com", "role": "driver" }
  }
}
```
Send the token on protected routes via the header:
```
Authorization: Bearer <token>
```

**Reserve** — `POST /api/v1/reservations`
```json
{ "zone_id": 5, "license_plate": "ABC-1234" }
```
</details>

---

## 🚦 Concurrency Handling — "EV Spot Bottleneck"

**The problem:** Suppose a zone has a capacity of 20, with 19 already active. At the exact same millisecond, two drivers send a request. Both read "19 active", both succeed → now there are 21 cars in a 20-slot zone. This is a **race condition**.

**The solution (in the repository layer):**
The entire check-and-create operation is wrapped in a **GORM Transaction**, and the zone row is held with a **`FOR UPDATE` lock**. As a result, the second request must wait to read that row until the first transaction finishes — making the operation atomic.

```go
db.Transaction(func(tx *gorm.DB) error {
    // 1️⃣ Lock the zone row — other transactions will wait
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)

    // 2️⃣ Check for a duplicate active license plate
    // 3️⃣ Count active reservations for this zone
    if int(active) >= zone.TotalCapacity {
        return ErrZoneFull        // → 409 Conflict
    }

    // 4️⃣ Safely create the reservation
    return tx.Create(res).Error
})
```

These two layers of protection (row lock + unique partial index) together ensure capacity can never overflow and duplicate plates are blocked.

---

## ⚙️ Local Setup

### Prerequisites
- Go 1.22+ (this project was built on 1.26)
- PostgreSQL (local or NeonDB/Supabase)
- (optional) [Air](https://github.com/air-verse/air) for hot-reload

### Steps
```bash
# 1. Clone
git clone https://github.com/arifuddincoder/spotsync.git
cd spotsync

# 2. Download dependencies
go mod download

# 3. Create .env (see the table below)
cp .env.example .env

# 4. Run
go run ./cmd/main.go
# or with hot-reload:
air
```

Once running: `http://localhost:8080`

### 🔑 Environment Variables (`.env`)
| Key | Description |
| --- | --- |
| `PORT` | Server port (e.g. `8080`) |
| `DSN` | PostgreSQL connection string |
| `JWT_SECRET_KEY` | Secret used to sign JWTs |
| `ADMIN_NAME` | Default admin name |
| `ADMIN_EMAIL` | Default admin email (used for seeding) |
| `ADMIN_PASSWORD` | Default admin password |

Example:
```env
PORT=8080
DSN="postgres://user:pass@localhost:5432/spotsync?sslmode=disable"
JWT_SECRET_KEY="your-super-secret-key"
ADMIN_NAME="Admin"
ADMIN_EMAIL="admin@spotsync.com"
ADMIN_PASSWORD="admin-password"
```

> 🌱 On server boot, an admin account is auto-seeded using these credentials (if one doesn't already exist).

---

## 🚀 Deployment

- **Backend:** [Render](https://render.com) — Go is natively supported.
- **Database:** NeonDB / Supabase (managed PostgreSQL).
- **CORS:** Open for all origins (configured in `server/http.go`).
- During deployment, set all of the environment variables above in the Render dashboard.

Live: **https://spotsync-b9vo.onrender.com**

---

## 📡 HTTP Status Codes

| Code | When |
| --- | --- |
| `200` | Successful GET/PUT/DELETE |
| `201` | New resource created (POST) |
| `400` | Validation / invalid input |
| `401` | Missing / invalid / expired token |
| `403` | Valid token but insufficient permission |
| `404` | Resource not found |
| `409` | Business conflict (zone full / duplicate plate) |
| `500` | Server / DB error |

---

## 🎤 Interview Video

> [https://drive.google.com/file/d/17KuOrfNjdlCQju92CTl4RwkeVS1xmHWw/view?usp=sharing](https://drive.google.com/file/d/17KuOrfNjdlCQju92CTl4RwkeVS1xmHWw/view?usp=sharing)

---

## 👤 Author

**Arif Uddin** · [GitHub @arifuddincoder](https://github.com/arifuddincoder)

---

> Built clean, concurrent & well-documented. 🚀
