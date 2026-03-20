# Backend - EASM Asset Management API

Go REST API backend for asset management and EASM (External Attack Surface Management) scan jobs.

## ðŸ“‹ Features

### Asset Management
- **CRUD Operations**: Create, read, update, delete assets
- **Batch Operations**: Batch create and batch delete multiple assets
- **Search & Filter**: Search by name, filter by type and status
- **Pagination**: Limit and offset based pagination
- **Statistics**: Get asset counts by type and status

### Scan Management
- **8 Active Scanners**:
  - IP Geolocation: Resolve IP location data
  - Port Scanning: TCP port enumeration (localhost/private IPs only)
  - SSL/TLS Inspection: Certificate chain analysis
  - Technology Detection: Web server/framework identification
  
- **5 Passive Scanners**:
  - DNS: A, MX, NS record lookup
  - WHOIS: Domain registrar information
  - Subdomain Enumeration: Common subdomain discovery
  - Certificate Transparency: CT log entries
  - ASN Lookup: Autonomous System information

- **Scan Orchestration**:
  - Support scan types: `dns`, `whois`, `subdomain`, `cert_trans`, `asn`, `ip`, `port`, `ssl`, `tech`, `all`
  - Asynchronous job tracking with status monitoring
  - Result storage and retrieval by scan type

### Security
- CORS middleware enabled
- Active scan safety rules: localhost/private IP only for port & SSL scans
- Public IP rejection for security-sensitive operations

## ðŸ”§ Requirements

- **Go**: 1.25 or later
- **PostgreSQL**: 12+
- **Environment**: Linux, macOS, or Windows with CGO support

## ðŸš€ Quick Start

### 1. Environment Setup

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Edit `.env`:
```env
DB_DSN=postgresql://user:password@localhost:5432/assets_db?sslmode=disable
PORT=8080
```

### 2. Database Setup

Ensure PostgreSQL is running and create the database:

```bash
psql -U postgres -c "CREATE DATABASE assets_db;"
```

The backend auto-creates required tables on startup:
- `assets`: Asset records (name, type, status, IP/domain)
- `scan_jobs`: Scan job tracking (asset_id, scan_type, status)
- `scan_results`: Scan result storage (result data by type)

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Server

```bash
go run ./cmd/server
```

Server starts on `http://localhost:8080`

### 5. Verify Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-03-20T10:30:00Z"
}
```

## ðŸ“¡ API Endpoints

### Asset Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/assets` | Create asset |
| GET | `/assets` | List assets (paginated) |
| GET | `/assets/stats` | Asset statistics |
| GET | `/assets/count` | Total asset count |
| GET | `/assets/search` | Search assets by name |
| DELETE | `/assets/{id}` | Delete asset |
| POST | `/assets/batch` | Batch create assets |
| DELETE | `/assets/batch` | Batch delete assets |

### Scan Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/assets/{id}/scan` | Start scan job |
| GET | `/scan-jobs/{id}` | Get scan job status |
| GET | `/scan-jobs/{id}/results` | Get scan results |
| GET | `/assets/{id}/scans` | List scans for asset |
| GET | `/assets/{id}/results` | List all results for asset |
| GET | `/assets/{id}/dns` | Get latest DNS scan result |
| GET | `/assets/{id}/whois` | Get latest WHOIS result |
| GET | `/assets/{id}/subdomains` | Get latest subdomain scan result |

## ðŸ§ª Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Specific Package Tests

```bash
go test -v ./internal/domain
go test -v ./internal/handler
go test -v ./internal/service
go test -v ./internal/scanner
```

### Run Single Test

```bash
go test -v ./internal/service -run TestAssetService_Create
```

### Generate Coverage Report File

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Packages

- **internal/domain**: Asset model validation
- **internal/handler**: HTTP endpoint testing via mock repository
- **internal/service**: Business logic with mock repository
- **internal/scanner**: Scanner tests (active + passive)

## ðŸ³ Docker Deployment

### Build Docker Image

```bash
docker build -t assets-api:latest .
```

### Run with Docker Compose

From the root directory:

```bash
docker-compose up
```

This starts:
- PostgreSQL container (postgres:16-alpine)
- API container (assets-api)

Verify:
```bash
curl http://localhost:8080/health
```

### Connect to Database in Container

```bash
docker exec -it assets_db psql -U postgres -d assets_db -c "SELECT COUNT(*) FROM assets;"
```

## ðŸ”’ Safety Rules

### Active Scan Restrictions

Port scans and SSL scans only work on:
- `127.0.0.1` - Localhost
- `10.0.0.0/8` - Private range
- `172.16.0.0/12` - Private range
- `192.168.0.0/16` - Private range

**Public IPs are rejected for safety.**

### Example - Safe Port Scan

```bash
curl -X POST http://localhost:8080/assets/{id}/scan \
  -H "Content-Type: application/json" \
  -d '{"scan_type": "port"}'

# Works if asset IP is 127.0.0.1 or 192.168.x.x
```

## ðŸ“Š Example Usage

### Create Asset

```bash
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api.example.com",
    "type": "domain"
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "api.example.com",
  "type": "domain",
  "status": "active",
  "created_at": "2024-03-20T10:30:00Z"
}
```

### Start Scan

```bash
curl -X POST http://localhost:8080/assets/{id}/scan \
  -H "Content-Type: application/json" \
  -d '{"scan_type": "dns"}'
```

### Get Scan Results

```bash
curl http://localhost:8080/scan-jobs/{job_id}/results
```

## ðŸ—ï¸ Project Structure

```
backend/
â”œâ”€â”€ cmd/
â”‚   â””â”€â”€ server/
â”‚       â””â”€â”€ main.go           # HTTP server entry point
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ config/               # Database config
â”‚   â”œâ”€â”€ domain/               # Domain models & validation
â”‚   â”œâ”€â”€ handler/              # HTTP request handlers
â”‚   â”œâ”€â”€ model/                # Request/response models
â”‚   â”œâ”€â”€ scanner/              # Scanner implementations (9 types)
â”‚   â”œâ”€â”€ service/              # Business logic (Asset/Scan services)
â”‚   â”œâ”€â”€ storage/              # Repository implementations
â”‚   â””â”€â”€ validator/            # Input validation
â”œâ”€â”€ tests/                    # All unit tests (separated)
â”‚   â”œâ”€â”€ domain/               # Model/domain tests
â”‚   â”œâ”€â”€ handler/              # Handler tests
â”‚   â”œâ”€â”€ service/              # Service tests
â”‚   â””â”€â”€ scanner/              # Scanner tests
â”œâ”€â”€ migrations/               # SQL migration files
â”œâ”€â”€ Dockerfile                # Docker image definition
â”œâ”€â”€ go.mod                    # Go module dependencies
â””â”€â”€ README.md                 # This file
```

## ðŸ” CI/CD

GitHub Actions runs on each push:

```yaml
Jobs:
  - Go Tests: go test -cover ./...
  - Gosec: Security scanning
  - Gitleaks: Secret detection
  - Trivy: Container vulnerability scan
```

View at: `.github/workflows/ci.yml`

## ðŸ“ Notes

- Server auto-applies migrations on startup
- CORS headers allow `*` origin (adjust in production)
- Scan jobs run asynchronously; check status via `/scan-jobs/{id}`
- Scan "all" type runs all passive scanners; partial failures allowed
- Results stored as JSON for extensibility

## ðŸ¤ Contributing

Run tests before committing:

```bash
go test ./...
go test -cover ./...
```

## ðŸ“„ License

Part of CMC Security Intern Program Day 3 assignments.

---

**For full project guide and Day 3 assignments 3-6, see [root README.md](../README.md)**

