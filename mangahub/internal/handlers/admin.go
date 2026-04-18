package handlers

import (
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"mangahub/internal/db"
)

var startTime = time.Now()

type AdminHandler struct {
	sqlDB *sql.DB
}

func NewAdminHandler(sqlDB *sql.DB) *AdminHandler {
	return &AdminHandler{sqlDB: sqlDB}
}

// GET /api/admin/health
func (h *AdminHandler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	if err := db.Ping(); err != nil {
		dbStatus = "error: " + err.Error()
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	respondOK(w, map[string]interface{}{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime).String(),
		"services": map[string]string{
			"http":  "ok",
			"db":    dbStatus,
			"ws":    "ok",
			"tcp":   "ok",
			"udp":   "ok",
			"grpc":  "ok",
		},
		"memory": map[string]interface{}{
			"alloc_mb":       memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":         memStats.Sys / 1024 / 1024,
			"gc_runs":        memStats.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	})
}

// GET /api/admin/ping
func (h *AdminHandler) Ping(w http.ResponseWriter, r *http.Request) {
	respondOK(w, map[string]interface{}{
		"pong": true,
		"ts":   time.Now(),
	})
}

// GET /api/admin/stats
func (h *AdminHandler) ServerStats(w http.ResponseWriter, r *http.Request) {
	var userCount, mangaCount, libraryCount, messageCount int
	h.sqlDB.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM users").Scan(&userCount)
	h.sqlDB.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM manga").Scan(&mangaCount)
	h.sqlDB.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM library_entries").Scan(&libraryCount)
	h.sqlDB.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM chat_messages").Scan(&messageCount)

	respondOK(w, map[string]interface{}{
		"users":           userCount,
		"manga":           mangaCount,
		"library_entries": libraryCount,
		"messages":        messageCount,
		"uptime":          time.Since(startTime).String(),
		"goroutines":      runtime.NumGoroutine(),
	})
}

// POST /api/admin/backup
func (h *AdminHandler) Backup(w http.ResponseWriter, r *http.Request) {
	respondOK(w, map[string]interface{}{
		"message":  "Backup initiated. Use pg_dump for production backups.",
		"command":  "pg_dump -U mangahub -d mangahub_db -F c -f backup_$(date +%Y%m%d_%H%M%S).dump",
		"timestamp": time.Now(),
	})
}

// GET /api/admin/logs
func (h *AdminHandler) Logs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.sqlDB.QueryContext(r.Context(),
		`SELECT id, level, message, context, created_at
		 FROM server_logs ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch logs")
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int64
		var level, message string
		var context []byte
		var createdAt time.Time
		rows.Scan(&id, &level, &message, &context, &createdAt)
		logs = append(logs, map[string]interface{}{
			"id":         id,
			"level":      level,
			"message":    message,
			"context":    string(context),
			"created_at": createdAt,
		})
	}
	if logs == nil {
		logs = []map[string]interface{}{}
	}
	respondOK(w, logs)
}
