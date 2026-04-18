package router

import (
	"database/sql"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"mangahub/internal/handlers"
	"mangahub/internal/middleware"
	"mangahub/internal/realtime"
	"mangahub/internal/repositories"
)

func New(db *sql.DB, hub *realtime.Hub) http.Handler {
	// Force MIME types for Windows
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimid.Recoverer)
	r.Use(chimid.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	// ── Repositories ──────────────────────────────────────────────────────
	userRepo   := repositories.NewUserRepository(db)
	mangaRepo  := repositories.NewMangaRepository(db)
	libRepo    := repositories.NewLibraryRepository(db)
	histRepo   := repositories.NewHistoryRepository(db)
	chatRepo   := repositories.NewChatRepository(db)
	notifRepo  := repositories.NewNotificationRepository(db)

	// ── Handlers ──────────────────────────────────────────────────────────
	authH    := handlers.NewAuthHandler(userRepo)
	mangaH   := handlers.NewMangaHandler(mangaRepo)
	libH     := handlers.NewLibraryHandler(libRepo, histRepo, mangaRepo)
	chatH    := handlers.NewChatHandler(chatRepo, hub)
	notifH   := handlers.NewNotificationHandler(notifRepo)
	adminH   := handlers.NewAdminHandler(db)
	wsH      := handlers.NewWSHandler(hub)

	// ── Static + Templates ────────────────────────────────────────────────
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// ── Web pages (server-side rendered) ─────────────────────────────────
	r.Get("/", serveTemplate("web/templates/pages/index.html"))
	r.Get("/login", serveTemplate("web/templates/pages/login.html"))
	r.Get("/register", serveTemplate("web/templates/pages/register.html"))
	r.Get("/dashboard", serveTemplate("web/templates/pages/dashboard.html"))
	r.Get("/browse", serveTemplate("web/templates/pages/browse.html"))
	r.Get("/manga/{id}", serveTemplate("web/templates/pages/manga_detail.html"))
	r.Get("/library", serveTemplate("web/templates/pages/library.html"))
	r.Get("/history", serveTemplate("web/templates/pages/history.html"))
	r.Get("/chat", serveTemplate("web/templates/pages/chat.html"))
	r.Get("/notifications", serveTemplate("web/templates/pages/notifications.html"))
	r.Get("/stats", serveTemplate("web/templates/pages/stats.html"))
	r.Get("/settings", serveTemplate("web/templates/pages/settings.html"))
	r.Get("/admin", serveTemplate("web/templates/pages/admin.html"))

	// ── WebSocket ─────────────────────────────────────────────────────────
	r.With(middleware.Auth).Get("/ws", wsH.Handle)

	// ── Auth API ──────────────────────────────────────────────────────────
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
		r.Post("/logout", authH.Logout)
		r.With(middleware.Auth).Get("/me", authH.Me)
		r.With(middleware.Auth).Put("/me", authH.UpdateProfile)
		r.With(middleware.Auth).Post("/change-password", authH.ChangePassword)
	})

	// ── Manga API ─────────────────────────────────────────────────────────
	r.Route("/api/manga", func(r chi.Router) {
		r.Get("/", mangaH.List)
		r.Get("/popular", mangaH.Popular)
		r.Get("/{id}", mangaH.GetByID)
	})

	// ── Library API (protected) ───────────────────────────────────────────
	r.Route("/api/library", func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/", libH.List)
		r.Post("/", libH.Add)
		r.Get("/stats", libH.Stats)
		r.Put("/{id}/progress", libH.UpdateProgress)
		r.Delete("/{id}", libH.Delete)
	})

	// ── History ───────────────────────────────────────────────────────────
	r.With(middleware.Auth).Get("/api/history", libH.History)

	// ── Chat API ─────────────────────────────────────────────────────────
	r.Route("/api/chat", func(r chi.Router) {
		r.Get("/rooms", chatH.GetRooms)
		r.With(middleware.Auth).Get("/rooms/{id}/messages", chatH.GetMessages)
		r.With(middleware.Auth).Post("/rooms/{id}/messages", chatH.SendMessage)
	})

	// ── Notifications ─────────────────────────────────────────────────────
	r.Route("/api/notifications", func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/", notifH.List)
		r.Put("/read-all", notifH.MarkAllRead)
	})

	// ── Admin ─────────────────────────────────────────────────────────────
	r.Route("/api/admin", func(r chi.Router) {
		r.Get("/ping", adminH.Ping)
		r.Get("/health", adminH.Health)
		r.With(middleware.Auth).Get("/stats", adminH.ServerStats)
		r.With(middleware.Auth).Get("/logs", adminH.Logs)
		r.With(middleware.Auth).Post("/backup", adminH.Backup)
	})

	return r
}

func serveTemplate(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
