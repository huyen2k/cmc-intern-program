package handler

import (
	"assets-api/internal/domain"
	"assets-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AssetHandler struct {
	service     *service.AssetService
	scanService *service.ScanService
}

func NewAssetHandler(s *service.AssetService, scanService *service.ScanService) *AssetHandler {
	return &AssetHandler{service: s, scanService: scanService}
}

func (h *AssetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var asset domain.Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(asset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *AssetHandler) DeleteByID(w http.ResponseWriter, r *http.Request, id string) {
	deleted, err := h.service.DeleteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "asset not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AssetHandler) StartScan(w http.ResponseWriter, r *http.Request, assetID string) {
	var req domain.StartScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.scanService.StartScan(assetID, req.ScanType)
	if err != nil {
		if service.IsNotFoundErr(err) {
			http.Error(w, "asset not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(job)
}

func (h *AssetHandler) GetScanJob(w http.ResponseWriter, r *http.Request, jobID string) {
	job, err := h.scanService.GetJob(jobID)
	if err != nil {
		if service.IsNotFoundErr(err) {
			http.Error(w, "scan job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func (h *AssetHandler) GetScanResults(w http.ResponseWriter, r *http.Request, jobID string) {
	resp, err := h.scanService.GetJobResults(jobID)
	if err != nil {
		if service.IsNotFoundErr(err) {
			http.Error(w, "scan job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *AssetHandler) ListAssetScans(w http.ResponseWriter, r *http.Request, assetID string) {
	scans, err := h.scanService.ListAssetScans(assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(scans)
}

func (h *AssetHandler) ListAssetResults(w http.ResponseWriter, r *http.Request, assetID string) {
	results, err := h.scanService.ListAssetResults(assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (h *AssetHandler) GetAssetDNS(w http.ResponseWriter, r *http.Request, assetID string) {
	h.getLatestResultByType(w, assetID, domain.ScanTypeDNS)
}

func (h *AssetHandler) GetAssetWhois(w http.ResponseWriter, r *http.Request, assetID string) {
	h.getLatestResultByType(w, assetID, domain.ScanTypeWhois)
}

func (h *AssetHandler) GetAssetSubdomains(w http.ResponseWriter, r *http.Request, assetID string) {
	h.getLatestResultByType(w, assetID, domain.ScanTypeSubdomain)
}

func (h *AssetHandler) getLatestResultByType(w http.ResponseWriter, assetID, scanType string) {
	result, err := h.scanService.GetLatestAssetResultByType(assetID, scanType)
	if err != nil {
		if service.IsNotFoundErr(err) {
			http.Error(w, "scan result not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// bài1
func (h *AssetHandler) Stats(w http.ResponseWriter, r *http.Request) {

	stats, err := h.service.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// bài 1
func (h *AssetHandler) Count(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	count, err := h.service.Count(t, status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"count": count,
		"filters": map[string]string{
			"type":   t,
			"status": status,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// bài2
func (h *AssetHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Assets []domain.Asset `json:"assets"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ids, err := h.service.BatchCreate(req.Assets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"created": len(ids),
		"ids":     ids,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// bài 3
func (h *AssetHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "ids parameter required", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsParam, ",")

	deleted, notFound, err := h.service.BatchDelete(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"deleted":   deleted,
		"not_found": notFound,
	}

	json.NewEncoder(w).Encode(response)
}

// bonus
func (h *AssetHandler) List(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	typeFilter := r.URL.Query().Get("type")
	statusFilter := r.URL.Query().Get("status")

	data, total, err := h.service.List(page, limit, typeFilter, statusFilter)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// bonus
func (h *AssetHandler) Search(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "q parameter required", http.StatusBadRequest)
		return
	}

	results, err := h.service.Search(query)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(results)
}
