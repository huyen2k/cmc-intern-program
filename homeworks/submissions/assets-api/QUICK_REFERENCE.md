# Quick Reference Guide

## ðŸ“ Folder Structure at a Glance

```
assets-api/
â”œâ”€â”€ backend/              â† Go REST API
â”‚   â”œâ”€â”€ cmd/server/       â† Server entry point (main.go)
â”‚   â”œâ”€â”€ internal/         â† Business logic
â”‚   â”‚   â”œâ”€â”€ handler/      â† HTTP handlers
â”‚   â”‚   â”œâ”€â”€ service/      â† Business services
â”‚   â”‚   â”œâ”€â”€ scanner/      â† 9 scanner types
â”‚   â”‚   â””â”€â”€ ...
â”‚   â”œâ”€â”€ Dockerfile
â”‚   â”œâ”€â”€ go.mod
â”‚   â””â”€â”€ README.md
â”‚
â”œâ”€â”€ frontend/             â† React + Vite dashboard
â”‚   â”œâ”€â”€ src/
â”‚   â”‚   â”œâ”€â”€ App.jsx       â† Main component
â”‚   â”‚   â””â”€â”€ styles.css    â† All styling
â”‚   â”œâ”€â”€ Dockerfile
â”‚   â”œâ”€â”€ nginx.conf
â”‚   â”œâ”€â”€ package.json
â”‚   â””â”€â”€ README.md
â”‚
â”œâ”€â”€ docker-compose.yml    â† Full stack config
â”œâ”€â”€ README.md             â† This project guide
â”œâ”€â”€ REORGANIZATION.md     â† What changed & why
â””â”€â”€ QUICK_REFERENCE.md    â† This file
```

## ðŸš€ Common Commands

### Start Everything (Fastest)
```bash
cd assets-api
docker-compose up
# Then open http://localhost:3000
```

### Backend Development
```bash
cd assets-api/backend

# Setup
cp .env.example .env
go mod tidy

# Run
go run ./cmd/server

# Test
go test ./...
go test -cover ./...
go test -v ./internal/handler
```

### Frontend Development
```bash
cd assets-api/frontend

# Setup
npm install

# Run
npm run dev
# Open http://localhost:5173

# Build
npm run build

# Lint
npm run lint
```

### Docker Commands
```bash
# Start services
docker-compose up

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild images
docker-compose build --no-cache

# Run specific service
docker-compose up postgres
docker-compose up backend
docker-compose up frontend
```

### Database Commands
```bash
# Connect to database in container
docker exec -it assets_db psql -U postgres -d assets

# Check tables
\dt

# Count assets
SELECT COUNT(*) FROM assets;

# Exit
\q
```

## ðŸ”— API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Health check |
| POST | `/assets` | Create asset |
| GET | `/assets` | List assets |
| GET | `/assets/stats` | Get statistics |
| GET | `/assets/search?q=name` | Search assets |
| DELETE | `/assets/{id}` | Delete asset |
| POST | `/assets/batch` | Batch create |
| DELETE | `/assets/batch` | Batch delete |
| POST | `/assets/{id}/scan` | Start scan |
| GET | `/scan-jobs/{id}` | Get job status |
| GET | `/scan-jobs/{id}/results` | Get results |
| GET | `/assets/{id}/dns` | Latest DNS scan |
| GET | `/assets/{id}/whois` | Latest WHOIS scan |
| GET | `/assets/{id}/subdomains` | Latest subdomain scan |

## ðŸ§ª Testing Quick Start

### Run All Backend Tests
```bash
cd backend
go test ./...
```

### Run Specific Tests
```bash
# Tests for a package
go test -v ./internal/domain
go test -v ./internal/handler
go test -v ./internal/service
go test -v ./internal/scanner

# Tests matching a pattern
go test -run TestAsset ./...

# With coverage
go test -cover ./...
```

### Frontend Linting
```bash
cd frontend
npm run lint
```

## ðŸ”„ Workflow Examples

### Add a New Scan Type

1. Create scanner file: `backend/internal/scanner/newscan_scanner.go`
2. Implement Scanner interface
3. Add to scanner map in `backend/internal/service/scan_service.go`
4. Add route in `backend/cmd/server/main.go`
5. Add tests in `backend/internal/scanner/newscan_scanner_test.go`
6. Test: `go test ./...`

### Update Dashboard UI

1. Edit `frontend/src/App.jsx` - Main component
2. Update `frontend/src/styles.css` - Styling
3. Check: `npm run lint`
4. Test: `npm run dev` then view http://localhost:5173

### Deploy to Production

```bash
# Build images
docker build -t assets-api:v1.0 -f backend/Dockerfile backend/
docker build -t assets-ui:v1.0 -f frontend/Dockerfile frontend/

# Push to registry
docker push yourregistry/assets-api:v1.0
docker push yourregistry/assets-ui:v1.0

# Deploy
docker-compose -f docker-compose.prod.yml up
```

## ðŸŒ Environment Variables

### Backend (backend/.env)
```env
DB_DSN=postgresql://user:pass@host:5432/assets
PORT=8080
```

### Frontend (frontend/.env.local)
```env
VITE_API_URL=http://localhost:8080
```

### Docker Compose
Already configured in `docker-compose.yml`:
- Backend: `DB_DSN=postgresql://postgres:postgres@postgres:5432/assets`
- Frontend: `VITE_API_URL=http://backend:8080`

## ðŸ“Š Project Stats

| Aspect | Details |
|--------|---------|
| **Language** | Go 1.25+ / React 19+ |
| **LOC** | ~2000 Go + ~3000 JS/CSS |
| **Endpoints** | 14+ API routes |
| **Scanners** | 9 types (4 active, 5 passive) |
| **Tests** | Comprehensive coverage 66-83% |
| **Containers** | 3 services (postgres, backend, frontend) |

## ðŸ› Troubleshooting

### Backend won't start
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check if port 8080 is free
lsof -i :8080

# Check environment
cd backend && cat .env
```

### Frontend won't connect to API
```bash
# Check backend health
curl http://localhost:8080/health

# Check .env.local
cd frontend && cat .env.local

# Check CORS
curl -X OPTIONS http://localhost:8080/assets
```

### Docker issues
```bash
# Remove old containers
docker-compose down -v

# Rebuild
docker-compose build --no-cache

# Check logs
docker-compose logs backend
docker-compose logs frontend
```

## ðŸ“š Documentation Files

| File | Purpose |
|------|---------|
| README.md | Complete project guide |
| REORGANIZATION.md | Migration guide & structure explanation |
| backend/README.md | Backend API documentation |
| frontend/README.md | Frontend dashboard guide |
| QUICK_REFERENCE.md | This file - quick commands |

## ðŸ”— File Navigation

```bash
# Go to backend from anywhere
cd $PROJECT_ROOT/backend

# Go to frontend from anywhere
cd $PROJECT_ROOT/frontend

# View specific files
cat backend/go.mod               # Backend dependencies
cat frontend/package.json        # Frontend dependencies
cat docker-compose.yml           # Docker configuration
```

## âœ¨ Key Files to Edit

| Task | File |
|------|------|
| Add new endpoint | `backend/cmd/server/main.go` |
| Add new scanner | `backend/internal/scanner/` |
| Fix business logic | `backend/internal/service/` |
| Update dashboard | `frontend/src/App.jsx` |
| Change styling | `frontend/src/styles.css` |
| Configure API URL | `frontend/.env.local` |
| Database schema | `backend/migrations/` |

## ðŸŽ¯ Day 3 Assignments

All core assignments documented in `README.md`:

- **BÃ i 1** âœ… Expand Scan API - 5 new scanners
- **BÃ i 2** âœ… Unit Tests - Full test coverage  
- **BÃ i 3** âœ… CORS Middleware - Frontend integration ready
- **BÃ i 4** âœ… CI/CD Pipeline - Automated testing
- **BÃ i 5** âœ… Docker Deployment - Full-stack containerization

## ðŸ’¡ Pro Tips

1. **Use docker-compose** for development - keeps dependencies isolated
2. **Check logs** when things fail: `docker-compose logs servicename`
3. **Environment variables** go in `.env` or `.env.local` - never commit
4. **Run tests** before committing: `go test ./...`
5. **Use Dockerfile** for reproducible builds
6. **CORS errors?** Check frontend `.env.local` has correct API URL
7. **Port conflicts?** Change in docker-compose.yml or use different ports

## ðŸš€ Next Steps

1. Try running: `docker-compose up`
2. Access dashboard: http://localhost:3000
3. Test creating an asset
4. Start a scan job
5. View results
6. Read `README.md` for full documentation
7. Read `REORGANIZATION.md` for what changed

---

**Happy coding!** ðŸŽ‰

