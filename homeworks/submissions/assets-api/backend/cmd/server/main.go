package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"assets-api/internal/config"
	"assets-api/internal/database"
	"assets-api/internal/handler"
	"assets-api/internal/repository"
	"assets-api/internal/service"
)

func isAPIPath(path string) bool {
	return strings.HasPrefix(path, "/assets") || strings.HasPrefix(path, "/scan-jobs") || path == "/health"
}

func spaStaticHandler(frontendDir string) http.HandlerFunc {
	indexFile := filepath.Join(frontendDir, "index.html")
	absFrontendDir, _ := filepath.Abs(frontendDir)

	return func(w http.ResponseWriter, r *http.Request) {
		if isAPIPath(r.URL.Path) {
			http.NotFound(w, r)
			return
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if cleanPath == "." {
			cleanPath = ""
		}

		requestedFile := filepath.Join(frontendDir, cleanPath)
		absRequestedFile, err := filepath.Abs(requestedFile)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		rel, err := filepath.Rel(absFrontendDir, absRequestedFile)
		if err != nil || strings.HasPrefix(rel, "..") {
			http.NotFound(w, r)
			return
		}

		if cleanPath != "" {
			if info, err := os.Stat(absRequestedFile); err == nil && !info.IsDir() {
				http.ServeFile(w, r, absRequestedFile)
				return
			}
		}

		if _, err := os.Stat(indexFile); err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, indexFile)
	}
}

func main() {

	cfg := config.Load()

	db, err := database.ConnectWithRetry(cfg.DB_DSN, 5)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewAssetRepository(db)
	assetService := service.NewAssetService(repo)
	scanRepo := repository.NewScanRepository(db)
	if err := scanRepo.EnsureSchema(); err != nil {
		log.Fatal(err)
	}
	scanService := service.NewScanService(assetService, scanRepo)

	assetHandler := handler.NewAssetHandler(assetService, scanService)

	mux := http.NewServeMux()

	mux.HandleFunc("/assets/stats", assetHandler.Stats)
	mux.HandleFunc("/assets/count", assetHandler.Count)

	mux.HandleFunc("/assets/batch", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			assetHandler.BatchCreate(w, r)
			return
		}

		if r.Method == http.MethodDelete {
			assetHandler.BatchDelete(w, r)
			return
		}
	})

	mux.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			assetHandler.List(w, r)
		case http.MethodPost:
			assetHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/assets/search", assetHandler.Search)
	mux.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/assets/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}

		id := parts[0]
		if len(parts) == 1 {
			if r.Method == http.MethodDelete {
				assetHandler.DeleteByID(w, r, id)
				return
			}
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		switch parts[1] {
		case "scan":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.StartScan(w, r, id)
		case "scans":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.ListAssetScans(w, r, id)
		case "results":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.ListAssetResults(w, r, id)
		case "dns":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.GetAssetDNS(w, r, id)
		case "whois":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.GetAssetWhois(w, r, id)
		case "subdomains":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.GetAssetSubdomains(w, r, id)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/scan-jobs/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/scan-jobs/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}

		jobID := parts[0]
		if len(parts) == 1 {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			assetHandler.GetScanJob(w, r, jobID)
			return
		}

		if parts[1] == "results" && r.Method == http.MethodGet {
			assetHandler.GetScanResults(w, r, jobID)
			return
		}

		http.NotFound(w, r)
	})

	// gọi đúng HealthHandler
	mux.HandleFunc("/health", handler.HealthHandler(db))

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "./web"
	}
	mux.HandleFunc("/", spaStaticHandler(frontendDir))

	log.Println("Server running on :" + cfg.PORT)

	http.ListenAndServe(":"+cfg.PORT, handler.CORSMiddleware(mux))
}
