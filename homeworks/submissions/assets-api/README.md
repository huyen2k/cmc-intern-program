# EASM Asset Management Platform

This is a full-stack Day 3 submission with a Go backend and React frontend.
This is the only README file in the project and includes all run instructions.

## Project Structure

```text
assets-api/
|-- backend/
|   |-- cmd/server/
|   |-- internal/
|   |-- migrations/
|   |-- go.mod
|   `-- Dockerfile
|-- frontend/
|   |-- src/
|   |-- package.json
|   `-- Dockerfile
|-- Dockerfile
`-- docker-compose.yml
```

## Day 3 Status

- Task 1: Expand Scan API - Done
- Task 2: Unit Tests - Done
- Task 3: Frontend Integration - Done
- Task 4: CI/CD - Done
- Task 5: Docker Compose Deployment - Done
- Task 6-9 (Bonus): Not implemented yet

## Run Project

### Option 1: Docker Compose (recommended)

Run from the assets-api folder:

```bash
docker-compose up --build
```

Access:

- App (frontend + API in one image): http://localhost:8080
- Health endpoint: http://localhost:8080/health

### Option 2: Local backend only

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

Backend runs on http://localhost:8080

### Option 3: Local frontend only

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on http://localhost:5173

## Run Tests

### Run all backend tests

```bash
cd backend
go test ./...
```

### Run tests with coverage

```bash
cd backend
go test -cover ./...
```

### Run specific package tests

```bash
cd backend
go test -v ./internal/domain
go test -v ./internal/handler
go test -v ./internal/service
go test -v ./internal/scanner
```

### Run one test by name

```bash
cd backend
go test -v ./internal/service -run TestAssetService_Create
```

### Generate HTML coverage report

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Docker Hub Image

Single-image deployment tag:

- huyennguyen08032005/assets-api:latest
