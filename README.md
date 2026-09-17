# Library & Reading Tracker REST API

A lightweight RESTful API for tracking books, reading progress, and personal reviews, built with **Go (Golang)** and **PostgreSQL**.

---

## Tech Stack

- **Language:** Go 1.26.2
- **Router:** [Gorilla Mux](https://github.com/gorilla/mux)
- **Database:** PostgreSQL (via [pgxpool](https://github.com/jackc/pgx))
- **Identifiers:** UUID v4 ([google/uuid](https://github.com/google/uuid))
- **Logging:** Structured logging with `log/slog` (pretty-printed console logs)
- **Migrations:** SQL migrations via [golang-migrate](https://github.com/golang-migrate/migrate)

---

## Features

- **Book Management:** Full CRUD operations for books (title, author, number of pages, year of publication).
- **Reading Tracker:** Mark books as read (`is_read`), tracking reading completion timestamps (`read_at`).
- **Reviews & Notes:** Store personal book reviews and impressions.
- **UUID Identifiers:** Primary keys use UUIDs instead of auto-incrementing integers.
- **Database Connection Pooling:** High-performance database operations using `pgxpool`.
- **Versioned Schema:** Database structure managed via sequential SQL migration files.

---

## Book Model

```json
{
  "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "number_of_pages": 464,
  "year_of_publication": 2008,
  "review": "A must-read for any software engineer.",
  "is_read": true,
  "read_at": "2026-03-25T14:30:00Z",
  "created_at": "2026-03-01T10:00:00Z"
}
```

---

## Project Structure

```text
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── api/                 # HTTP server, handlers, DTOs & middleware
│   ├── db/                  # PostgreSQL connection pool setup (pgxpool)
│   ├── library/             # Domain models, repository & storage interfaces
│   └── logger/              # slog setup & slogpretty handler
├── migrations/              # SQL migration files (up/down)
└── go.mod                   # Go module dependencies
```

---

## Getting Started

### Prerequisites
- [Go](https://go.dev/) (version 1.26)
- [PostgreSQL](https://www.postgresql.org/) database running locally

### Installation & Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/loundxr/library-rest-api.git
   cd library-rest-api
   ```

2. **Apply Migrations:**
   Run the SQL scripts located in the `migrations/` directory against your PostgreSQL database.

3. **Run the server:**
   ```bash
   go run cmd/main.go
   ```