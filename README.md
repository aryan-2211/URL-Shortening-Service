# URL-Shortening-Service

A URL shortening service written in Go. Takes a long URL, stores it in a local SQLite database, and returns a 6-character short code. Visiting the short URL redirects you to the original.

## How It Works

The service exposes two HTTP endpoints:

| Endpoint | Method | Description |
|---|---|---|
| `/create?longURL=<url>` | GET | Generates a short code for the given URL and persists it |
| `/<shortCode>` | GET | Looks up the short code and redirects to the original URL |

Short codes are 6 characters long, randomly sampled from `[a-zA-Z0-9]`, giving 62⁶ (~56 billion) possible values. The mapping is stored in a local `urls.db` SQLite file using GORM with `AutoMigrate` — no manual schema setup needed.

```
POST /create?longURL=https://example.com/very/long/path
        │
        ▼
  generateShortURL()    ← 6-char random alphanumeric
        │
        ▼
  db.Create(&URL{})     ← persisted to SQLite via GORM
        │
        ▼
  Response: "Short URL: aB3xYz"


GET /aB3xYz
        │
        ▼
  db.First(&url, "short_url = ?", "aB3xYz")
        │
        ▼
  http.Redirect → https://example.com/very/long/path  (302)
```

## Project Structure

```
URL-Shortening-Service/
├── main.go      # All logic: DB setup, handlers, short code generation
├── go.mod
└── go.sum
```

## Data Model

```go
type URL struct {
    ID       uint   `gorm:"primaryKey"`
    LongURL  string `gorm:"not null"`
    ShortURL string `gorm:"uniqueIndex"`
}
```

The `uniqueIndex` on `ShortURL` ensures no two entries share the same short code. The database file (`urls.db`) is created automatically on first run.

## Running Locally

**Prerequisites:** Go 1.19+

```bash
git clone https://github.com/aryan-2211/URL-Shortening-Service.git
cd URL-Shortening-Service
go mod tidy
go run main.go
```

The server starts on `http://localhost:8080`.

## Usage

**Create a short URL:**
```bash
curl "http://localhost:8080/create?longURL=https://www.github.com/aryan-2211"
# Short URL: aB3xYz
```

**Use the short URL** — open in browser or curl with redirect following:
```bash
curl -L "http://localhost:8080/aB3xYz"
# → 302 redirect to https://www.github.com/aryan-2211
```

**Missing `longURL` parameter:**
```bash
curl "http://localhost:8080/create"
# 400 Bad Request: LongURL is required
```

**Unknown short code:**
```bash
curl "http://localhost:8080/notfound"
# 404 Not Found
```

## Tech Stack

- **Language:** Go 1.19
- **Database:** SQLite (via `glebarez/sqlite`, CGO-free driver)
- **ORM:** GORM v1.25 with `AutoMigrate`
- **HTTP:** Go standard library `net/http`
