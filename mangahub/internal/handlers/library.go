package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mangahub/internal/middleware"
	"mangahub/internal/models"
	"mangahub/internal/repositories"
)

type LibraryHandler struct {
	library *repositories.LibraryRepository
	history *repositories.HistoryRepository
	manga   *repositories.MangaRepository
}

func NewLibraryHandler(
	library *repositories.LibraryRepository,
	history *repositories.HistoryRepository,
	manga   *repositories.MangaRepository,
) *LibraryHandler {
	return &LibraryHandler{library: library, history: history, manga: manga}
}

// GET /api/library?status=
func (h *LibraryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	status := r.URL.Query().Get("status")

	entries, err := h.library.GetByUser(r.Context(), uid, status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch library")
		return
	}
	if entries == nil {
		entries = []*models.LibraryEntry{}
	}
	respondOK(w, entries)
}

// POST /api/library
func (h *LibraryHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)

	var req models.AddLibraryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.MangaID == uuid.Nil {
		respondError(w, http.StatusBadRequest, "manga_id is required")
		return
	}

	// Verify manga exists
	m, _ := h.manga.GetByID(r.Context(), req.MangaID)
	if m == nil {
		respondError(w, http.StatusNotFound, "manga not found")
		return
	}

	now := time.Now()
	entry := &models.LibraryEntry{
		UserID:  uid,
		MangaID: req.MangaID,
		Status:  req.Status,
	}
	if req.Status == "reading" {
		entry.StartedAt = &now
	}

	if err := h.library.Add(r.Context(), entry); err != nil {
		respondError(w, http.StatusConflict, "manga already in library")
		return
	}
	respondCreated(w, entry)
}

// PUT /api/library/{id}/progress
func (h *LibraryHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entry ID")
		return
	}

	entry, _ := h.library.GetByID(r.Context(), entryID, uid)
	if entry == nil {
		respondError(w, http.StatusNotFound, "library entry not found")
		return
	}

	var req models.UpdateProgressRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate chapter progress
	m, _ := h.manga.GetByID(r.Context(), entry.MangaID)
	if m != nil && m.ChapterCount > 0 && req.CurrentChapter > m.ChapterCount {
		respondError(w, http.StatusBadRequest, "chapter exceeds total chapter count")
		return
	}

	oldChapter := entry.CurrentChapter
	entry.CurrentChapter = req.CurrentChapter
	entry.CurrentVolume = req.CurrentVolume
	if req.Status != "" {
		entry.Status = req.Status
	}
	entry.Rating = req.Rating
	entry.Notes = req.Notes

	if err := h.library.UpdateProgress(r.Context(), entry); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update progress")
		return
	}

	// Record history if chapter advanced
	if req.CurrentChapter > oldChapter {
		h.history.Add(r.Context(), &models.ReadingHistory{
			UserID:        uid,
			MangaID:       entry.MangaID,
			ChapterNumber: req.CurrentChapter,
			VolumeNumber:  req.CurrentVolume,
		})
	}

	respondOK(w, entry)
}

// DELETE /api/library/{id}
func (h *LibraryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	entryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entry ID")
		return
	}
	if err := h.library.Delete(r.Context(), entryID, uid); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete entry")
		return
	}
	respondOK(w, map[string]string{"message": "removed from library"})
}

// GET /api/library/stats
func (h *LibraryHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	stats, err := h.library.GetStats(r.Context(), uid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}
	respondOK(w, stats)
}

// GET /api/history
func (h *LibraryHandler) History(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	history, err := h.history.GetByUser(r.Context(), uid, 50, 0)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}
	if history == nil {
		history = []*models.ReadingHistory{}
	}
	respondOK(w, history)
}
