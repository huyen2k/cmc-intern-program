# Backend - Asset Management API

## Quick Start

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

- `assets`: asset records
- `scan_jobs`: scan job tracking
- `scan_results`: scan result storage

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Server

```bash
go run ./cmd/server
```

Server starts on `http://localhost:8080`.

### 5. Verify Health

```bash
curl http://localhost:8080/health
```

## Testing

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

## Docker Deployment

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

- PostgreSQL container
- Backend API container

## Safety Rules

Port and SSL scans are restricted to private or localhost ranges.
Public IPs are blocked for active scans.

## Example API Usage

### Create Asset

```bash
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api.example.com",
    "type": "domain"
  }'
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

## Project Structure

```text
backend/
|-- cmd/
|   `-- server/
|       `-- main.go
|-- internal/
|   |-- config/
|   |-- domain/
|   |-- handler/
|   |-- model/
|   |-- scanner/
|   |-- service/
|   |-- storage/
|   `-- validator/
|-- migrations/
|-- Dockerfile
|-- go.mod
`-- README.md
```

## CI/CD

GitHub Actions workflow: `.github/workflows/ci.yml`

## Notes

- Migrations are applied on startup.
- Scan jobs run asynchronously.
- Results are stored as JSON.

## License

Part of CMC Security Intern Program Day 3 assignments.
