package handler

import (
	"encoding/json"
	"link-storage-service/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	service service.LinkService
}

func NewLinkHandler(s service.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

type createRequest struct {
	URL string `json:"url"`
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	code, err := h.service.Create(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"short_code": code,
	})
}

func (h *LinkHandler) GetByShortCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.Get(code)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"url":    link.OriginalURL,
		"visits": link.Visits,
	})
}

func (h *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	links, err := h.service.GetAll(limit, offset)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) DeleteLinks(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	err := h.service.Delete(code)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetStats(code)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"short_code": link.ShortCode,
		"url":        link.OriginalURL,
		"visits":     link.Visits,
		"created_at": link.CreatedAt,
	})
}
