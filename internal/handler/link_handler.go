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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	code, err := h.service.Create(req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"short_code": code,
	})
}

func (h *LinkHandler) GetByShortCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.Get(code)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"url":    link.OriginalURL,
		"visits": link.Visits,
	})
}

func (h *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	links, err := h.service.GetAll(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error")
		return
	}

	writeJSON(w, http.StatusOK, links)
}

func (h *LinkHandler) DeleteLinks(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if err := h.service.Delete(code); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	link, err := h.service.GetStats(code)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"short_code": link.ShortCode,
		"url":        link.OriginalURL,
		"visits":     link.Visits,
		"created_at": link.CreatedAt,
	})
}
