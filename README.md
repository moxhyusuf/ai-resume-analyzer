# AI Resume Analyzer API

REST API for AI-powered resume analysis, built with **Go + Fiber** and backed by **Groq** (`openai/gpt-oss-20b`).

## Stack

| Layer | Tech |
|---|---|
| Language | Go 1.26.3 |
| HTTP Framework | Fiber v2 |
| Database | PostgreSQL 17 + GORM |
| AI Provider | Groq API — `openai/gpt-oss-20b` |
| Auth | JWT |
| Dev Hot-reload | Air + Swagger (swag) |
| Containerization | Docker + Docker Compose |

## Why Groq + gpt-oss-20b?

- **1000+ tokens/sec** — Groq's LPU hardware, near-instant analysis
- OpenAI-compatible REST API — no SDK needed, plain HTTP
- Apache 2.0 licensed — free for commercial use
- MoE architecture: 20B params, only 3.6B active per token (low latency)
- Pricing: $0.10/M input · $0.50/M output

---

## Project Structure

```
ai-resume-analyzer/
├── cmd/api/main.go               # Entry point, dependency wiring
├── config/config.go              # Env loading, DB connect, AutoMigrate
├── internal/
│   ├── auth/
│   │   ├── dto.go
│   │   ├── handler.go
│   │   └── service.go
│   ├── middleware/
│   │   └── jwt.go
│   ├── resume/
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go            # PDF extract → Groq AI → save
│   └── user/
│       ├── model.go
│       └── repository.go
├── pkg/
│   ├── groq/client.go            # Groq REST client (OpenAI-compatible)
│   ├── pdf/extractor.go          # PDF text extraction
│   ├── response/response.go      # Standardized JSON response
│   └── validator/validator.go
├── docs/                         # Swagger (auto-generated)
├── postman/                      # Postman collection
├── Dockerfile                    # Production build (distroless)
├── Dockerfile.dev                # Dev build with Air + swag
├── docker-compose.yml            # Production compose
├── docker-compose.dev.yml        # Dev compose (hot-reload + exposed DB port)
├── .air.toml                     # Air config (rebuilds + swag regen on save)
├── .env.example
└── Makefile
```

---

## Getting Started

### Local (without Docker)

```bash
go mod tidy
cp .env.example .env
# Edit .env — set GROQ_API_KEY from console.groq.com
make run
```

### Dev with Hot-reload (Air)

```bash
make dev
# or
air -c .air.toml
```

Air will rebuild on every `.go` / `.env` change and regenerate Swagger docs automatically via `swag init`.

### Docker — Development

```bash
make docker-dev-up
# or
docker compose -f docker-compose.dev.yml up --build
```

- App + PostgreSQL in containers
- Source mounted into container (`./:/app`) — hot-reload active
- DB exposed on `localhost:5432` for direct access (e.g. TablePlus, DBeaver)

### Docker — Production

```bash
make docker-prod
# or
docker compose up --build
```

- Multi-stage build: `golang:1.26.3-alpine` → `distroless/static-debian12:nonroot`
- App port bound to `127.0.0.1:8080` (not public-facing — put Nginx/Caddy in front)
- DB **not** exposed to host

---

## Environment Variables

```env
PORT=8080
DB_USER=resume_user
DB_PASSWORD=secret
DB_NAME=ai_resume_analyzer
DATABASE_URL=host=postgres user=resume_user password=secret dbname=ai_resume_analyzer port=5432 sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxxxxxx
```

---

## API Endpoints

### Auth

```
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Resumes — Bearer token required

```
POST   /api/v1/resumes/analyze    multipart/form-data, field: resume (PDF, ≤5MB)
GET    /api/v1/resumes            list all user's resumes
GET    /api/v1/resumes/:id        detail + analysis result
DELETE /api/v1/resumes/:id        delete resume
```

### Analyze Response Example

```json
{
  "success": true,
  "message": "Resume analyzed successfully",
  "data": {
    "resume_id": "uuid-here",
    "analysis": {
      "overall_score": 78,
      "summary_feedback": "Strong experience section but skills need more specifics.",
      "sections": [
        { "name": "Experience", "score": 85, "feedback": "Good use of action verbs." },
        { "name": "Skills",     "score": 60, "feedback": "Add proficiency levels." },
        { "name": "Education",  "score": 80, "feedback": "Clear and concise." }
      ],
      "suggestions": [
        "Add quantifiable achievements (e.g. improved API response time by 40%)",
        "Include a professional summary at the top"
      ],
      "keywords": ["Go", "REST API", "PostgreSQL", "Docker"]
    }
  }
}
```

Swagger UI availabe at `http://localhost:8080/swagger/`

Postman collection availabe at `postman/`

---

## AI Integration Flow

```
POST /resumes/analyze
  ↓
handler       → validate file (PDF, ≤5MB)
  ↓
pdf.Extract   → extract plain text from PDF bytes
  ↓
groq.Client   → POST api.groq.com/openai/v1/chat/completions
                 model: openai/gpt-oss-20b
                 system: structured JSON resume reviewer prompt
  ↓
service       → json.Unmarshal AI response into AnalysisResult struct
  ↓
repository    → save Resume + AnalysisResult to PostgreSQL
  ↓
handler       → return standardized JSON response
```

---

## Makefile Commands

| Command | Description |
|---|---|
| `make run` | `go run ./cmd/api` |
| `make build` | Build binary to `bin/api` |
| `make dev` | Hot-reload via Air |
| `make tidy` | `go mod tidy && go mod verify` |
| `make lint` | Run `golangci-lint` |
| `make clean` | Remove `bin/` and `tmp/` |
| `make docker-dev-up` | Start dev compose (with build) |
| `make docker-dev-down` | Stop dev compose |
| `make docker-prod` | Start production compose (with build) |