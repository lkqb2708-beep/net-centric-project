package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mangahub/internal/repositories"
)

type MangaHandler struct {
	manga *repositories.MangaRepository
}

func NewMangaHandler(manga *repositories.MangaRepository) *MangaHandler {
	return &MangaHandler{manga: manga}
}

// GET /api/manga?q=&status=&genre=&page=&page_size=
func (h *MangaHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("q")
	status := q.Get("status")
	genre := q.Get("genre")
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	manga, total, err := h.manga.Search(r.Context(), query, status, genre, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch manga")
		return
	}
	totalPages := (total + pageSize - 1) / pageSize
	respondOK(w, map[string]interface{}{
		"data":        manga,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GET /api/manga/popular
func (h *MangaHandler) Popular(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 12
	}
	manga, err := h.manga.GetPopular(r.Context(), limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch popular manga")
		return
	}
	respondOK(w, manga)
}

// GET /api/manga/{id}
func (h *MangaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid manga ID")
		return
	}
	m, err := h.manga.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch manga")
		return
	}
	if m == nil {
		respondError(w, http.StatusNotFound, "manga not found")
		return
	}
	respondOK(w, m)
}
