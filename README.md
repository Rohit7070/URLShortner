# URL Shortener Service in Go

A simple URL shortener service built with **Go**, **Gin**, and **GORM**, supporting:

- Shortening long URLs
- Custom alias support
- Redirect via short code
- Viewing statistics for each short URL
- Rate limiting for API requests

---

## 📦 Features

1. **Shorten long URLs**
2. **Custom alias support**
3. **Redirect short URL to original URL**
4. **View statistics (hits, creation time)**
5. **Rate limiting** to prevent abuse

---

## 🗂️ Project Structure

URLShortnerAssignment/
├── cmd/
│ └── main.go # Entry point
├── internal/
│ ├── controller/ # Gin HTTP handlers
│ ├── service/ # Business logic
│ ├── repository/ # DB access (GORM)
│ ├── models/ # DB models
│ ├── interfaces/ # Service and repository interfaces
│ ├── middleware/ # Rate limiting middleware
│ └── test/ # Integration/API tests
├── go.mod
└── README.md



---

## ⚙️ Prerequisites

- Go >= 1.20
- SQLite (default, can be replaced by Postgres/MySQL)
- Git

---

## 🏃‍♂️ Running the Project

1. **Clone the repository**

```bash
git clone <https://github.com/Rohit7070/URLShortner>
cd URLShortnerAssignment

Install dependencies
go mod tidy

go run cmd/main.go

Server runs by default on:

http://localhost:8090

🔗 API Endpoints
1. Shorten URL
POST /shorten


Request Body:

{
  "long_url": "https://www.example.com/very/long/url",
  "custom": "chatgpt"      // optional
}


Responses:

201 Created — URL shortened successfully

409 Conflict — custom alias already exists

400 Bad Request — invalid request

429 Too Many Requests — rate limit exceeded

Example cURL:

curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"long_url":"https://www.example.com","custom":"chatgpt"}'

2. Redirect
GET /:code


Behavior:

Redirects to the original long URL

404 Not Found if code does not exist

Example cURL:

curl -v http://localhost:8080/chatgpt

3. Stats
GET /stats/:code


Response:

{
  "ID": 1,
  "ShortCode": "chatgpt",
  "LongURL": "https://www.example.com",
  "Hits": 5,
  "CreatedAt": "2025-10-08T13:30:00Z"
}


404 Not Found if the code does not exist

Example cURL:

curl http://localhost:8080/stats/chatgpt

🧪 Running Tests

Integration tests are located in internal/test:

go test ./internal/test -v

⚡ Notes

Rate limiting is applied per IP for /shorten.

Auto-generated short codes are 6-character alphanumeric by default.

Custom aliases must be alphanumeric.

GORM handles database operations; errors are mapped to HTTP status codes.


This is a **complete README.md** ready to place in your project root.  

Do you want me to also **add a Quick Start section with Postman collection examples** for testing all APIs visually?




