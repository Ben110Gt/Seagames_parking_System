Motorcycle ticket payment system
Currently, parking at the Seagames dormitory is still using a day-based system, so we decided to build this system

📁 Requirements
1) A login system that separates the roles of the parking attendant and the owner
2) When the motorcycle enters, a ticket will be created and a ticket will be issued (there will be a barcode on the ticket). After the motorcycle leaves, the barcode will be scanned
3) If it exceeds 1 day, there will be a fine of 2000kip per day (for example, in the ticket on 12/4/2026, but when the barcode is scanned on 15/4/2026, there will be a fine of 6000kip)

4) There is a monthly membership system of 60,000kip
5) View daily, weekly, and monthly income
6) Can search for customer motorcycle tickets (in case the customer loses the ticket)

📁 Tech Stack

Language: Go (Golang) , Dart

Framework: Fiber v2 ,  Flutter

Database: PostgreSQL

ORM: GORM



Configuration: Environment-based configuration management.

## 📁 Project Structure

```
.
├── frontend/
│
└── backend/
    ├── cmd/
    │   └── api/
    │       └── main.go
    │
    ├── internal/
    │   ├── models/       # GORM model structs
    │   ├── repository/   # Database access layer
    │   ├── service/      # Business logic
    │   ├── handler/      # HTTP handlers (Fiber)
    │   ├── routes/       # Route registration
    │   ├── middleware/   # Auth, logging middleware
    │   ├── server/       # Fiber app setup
    │   └── utils/        # Helpers (QR, date calc, etc.)
    │
    └── database/         # Migration & DB connection
```

## Prerequisites

- **Go** (version specified in `go.mod`)
- **PostgreSQL** running locally or on a remote host

Verify your Go installation:

```bash
go version
```

---

## Environment Configuration

Create a `.env` file in the project root (the same directory as `go.mod`) with the following variables:

```env
# Application
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=seagame_parking
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Vientiane

# JWT
JWT_SECRET=your_jwt_secret_key

# Fine rate (optional — defaults to 2000 KIP/day if not set)
DAILY_FINE_KIP=2000
```

> **Note:** The database schema is created automatically via GORM `AutoMigrate` on first startup — no manual SQL migrations are required.

---

## Getting Started

**1. Clone the repository**

```bash
git clone https://github.com/your-username/seagame-parking-backend.git
cd seagame-parking-backend
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Set up the `.env` file** (see section above)

**4. Run the server**

```bash
go run ./cmd/api
```

The API will be available at `http://0.0.0.0:<APP_PORT>` (e.g. `http://localhost:8080`).

A health-check endpoint is available at:

```
GET /  →  { "status": "ok", "message": "Seagames Parking System API" }
```

---


---

## Authentication & Roles

All protected routes require a JWT token passed in the `Authorization` header:

```
Authorization: Bearer <token>
```

Tokens can also be passed as a `?token=` query parameter (useful for WebSocket connections).

| Role | Description |
|---|---|
| *(none)* | Public — no token required |
| `attendant` | Default role assigned on registration |
| `Owner` | Full administrative access |

JWT tokens expire after **24 hours**.

---

## API Endpoints

Base URL: `http://localhost:8080`

---

### Auth Routes `(public)`

| Method | Path | Description |
|---|---|---|
| `POST` | `/auth/register` | Register a new user account |
| `POST` | `/auth/login` | Log in and receive a JWT token |

---

### User Management `(role: Owner)`

| Method | Path | Description |
|---|---|---|
| `GET` | `/admin/users/` | List all users |
| `GET` | `/admin/users/:id` | Get a single user by ID |
| `PUT` | `/admin/users/:id` | Update a user |
| `DELETE` | `/admin/users/:id` | Delete a user |

---

### Ticket Management `(role: Owner)`

| Method | Path | Description |
|---|---|---|
| `POST` | `/admin/create-ticket` | Issue a new parking ticket (returns QR code) |
| `POST` | `/admin/check-ticket` | Check out a vehicle and calculate the fee |
| `GET` | `/admin/tickets/search` | Search tickets by plate number, code, or status |
| `GET` | `/admin/tickets/active` | List all currently active (parked) tickets |

#### Ticket & Fine Logic

- On **check-in**, a unique ticket code is generated and a QR code (base64 PNG) is returned.
- On **check-out**, the system calculates the fee:
  - **Base fee**: 2,000 KIP (flat)
  - **Daily fine**: 2,000 KIP per extra day (configurable via `DAILY_FINE_KIP`)
  - Same-day check-out incurs only the base fee — no penalty.
  - If the vehicle's plate number has an active membership, **no penalty** is applied.
- Membership QR codes (prefixed `MBR-`) can also be scanned at check-out — access is granted at no charge if the card is still valid.
- Checked-out tickets are **soft-deleted** (archived) to preserve the transaction record.

---

### Membership Management `(role: Owner)`

| Method | Path | Description |
|---|---|---|
| `POST` | `/admin/create-membership` | Create a new membership card |
| `GET` | `/admin/memberships` | List all membership cards |
| `GET` | `/admin/memberships/:code` | Get a membership card by code |
| `GET` | `/admin/memberships/active` | List active (non-expired) memberships *(in progress)* |

Each membership card is linked to a specific **plate number** and has a **registration date** and **expiration date**. The default fee is 60,000 KIP.

---

### Income Reports `(role: Owner)`

| Method | Path | Description |
|---|---|---|
| `GET` | `/admin/income/daily` | Income report for today |
| `GET` | `/admin/income/weekly` | Income report for the current week |
| `GET` | `/admin/income/monthly` | Income report for the current month |

Income is calculated from the `Transaction` table, which records both ticket fines and membership fees.

---



## Author

- Name: Souksakhone Haknolath
- Contact: souksakhone.haknolth@gmail.com
