# EASM Asset Management Platform

This repository contains a full-stack Day 3 submission:

- Backend API (Go): `backend/`
- Frontend Dashboard (React + Vite): `frontend/`

## Project Structure

```text
assets-api/
|-- backend/
|   |-- cmd/server/
|   |-- internal/
|   |-- migrations/
|   |-- go.mod
|   `-- README.md
|-- frontend/
|   |-- src/
|   |-- package.json
|   `-- README.md
`-- docker-compose.yml
```

## Day 3 Status

- Task 1: Expand Scan API - Done
- Task 2: Unit Tests - Done
- Task 3: Frontend Integration - Done
- Task 4: CI/CD - Done
- Task 5: Docker Compose Deployment - Done
- Task 6-9 (Bonus): Not implemented yet

## Run Tests

Run all backend tests:

```bash
cd backend
go test ./...
```

Run with coverage:

```bash
cd backend
go test -cover ./...
```

Run specific package tests:

```bash
cd backend
go test -v ./internal/domain
go test -v ./internal/handler
go test -v ./internal/service
go test -v ./internal/scanner
```

Run one specific test by name:

```bash
cd backend
go test -v ./internal/service -run TestAssetService_Create
```

Generate coverage report file:

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Quick Start

Docker (recommended):

```bash
docker-compose up --build
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080

## Documentation

- Backend details: `backend/README.md`
- Frontend details: `frontend/README.md`
