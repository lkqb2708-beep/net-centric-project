package models

import (
	"time"

	"github.com/google/uuid"
)

// ─── User ──────────────────────────────────────────────────────────────────

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	AvatarURL    string     `json:"avatar_url" db:"avatar_url"`
	Bio          string     `json:"bio" db:"bio"`
	Role         string     `json:"role" db:"role"` // user | admin
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastSeenAt   *time.Time `json:"last_seen_at" db:"last_seen_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type UserPublic struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	AvatarURL string     `json:"avatar_url"`
	Bio       string     `json:"bio"`
	Role      string     `json:"role"`
	LastSeenAt *time.Time `json:"last_seen_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// ─── Session ───────────────────────────────────────────────────────────────

type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ─── Manga ─────────────────────────────────────────────────────────────────

type Manga struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Author       string    `json:"author" db:"author"`
	Artist       string    `json:"artist" db:"artist"`
	Genres       []string  `json:"genres" db:"genres"`
	Status       string    `json:"status" db:"status"` // ongoing | completed | hiatus | cancelled
	ChapterCount int       `json:"chapter_count" db:"chapter_count"`
	VolumeCount  int       `json:"volume_count" db:"volume_count"`
	Description  string    `json:"description" db:"description"`
	CoverURL     string    `json:"cover_url" db:"cover_url"`
	Year         int       `json:"year" db:"year"`
	Rating       float64   `json:"rating" db:"rating"`
	PopularityRank int     `json:"popularity_rank" db:"popularity_rank"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ─── Library Entry ─────────────────────────────────────────────────────────

type LibraryEntry struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	MangaID        uuid.UUID  `json:"manga_id" db:"manga_id"`
	Status         string     `json:"status" db:"status"` // reading|completed|plan_to_read|on_hold|dropped
	CurrentChapter int        `json:"current_chapter" db:"current_chapter"`
	CurrentVolume  int        `json:"current_volume" db:"current_volume"`
	Rating         *float64   `json:"rating" db:"rating"`
	Notes          string     `json:"notes" db:"notes"`
	StartedAt      *time.Time `json:"started_at" db:"started_at"`
	FinishedAt     *time.Time `json:"finished_at" db:"finished_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`

	// Joined fields
	Manga *Manga `json:"manga,omitempty" db:"-"`
}

// ─── Reading History ────────────────────────────────────────────────────────

type ReadingHistory struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	MangaID       uuid.UUID `json:"manga_id" db:"manga_id"`
	ChapterNumber int       `json:"chapter_number" db:"chapter_number"`
	VolumeNumber  int       `json:"volume_number" db:"volume_number"`
	ReadAt        time.Time `json:"read_at" db:"read_at"`

	// Joined fields
	Manga *Manga `json:"manga,omitempty" db:"-"`
}

// ─── Review ────────────────────────────────────────────────────────────────

type Review struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	MangaID   uuid.UUID `json:"manga_id" db:"manga_id"`
	Rating    float64   `json:"rating" db:"rating"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	User  *UserPublic `json:"user,omitempty" db:"-"`
	Manga *Manga      `json:"manga,omitempty" db:"-"`
}

// ─── Friend ────────────────────────────────────────────────────────────────

type Friend struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FriendID  uuid.UUID `json:"friend_id" db:"friend_id"`
	Status    string    `json:"status" db:"status"` // pending | accepted | blocked
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Joined fields
	Friend *UserPublic `json:"friend,omitempty" db:"-"`
}

// ─── Activity Feed ─────────────────────────────────────────────────────────

type ActivityFeedItem struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	UserID     uuid.UUID              `json:"user_id" db:"user_id"`
	ActionType string                 `json:"action_type" db:"action_type"`
	MangaID    *uuid.UUID             `json:"manga_id" db:"manga_id"`
	Metadata   map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`

	// Joined fields
	User  *UserPublic `json:"user,omitempty" db:"-"`
	Manga *Manga      `json:"manga,omitempty" db:"-"`
}

// ─── Chat ──────────────────────────────────────────────────────────────────

type ChatRoom struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Type        string     `json:"type" db:"type"` // general | manga
	MangaID     *uuid.UUID `json:"manga_id" db:"manga_id"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`

	// Joined
	Manga        *Manga `json:"manga,omitempty" db:"-"`
	MemberCount  int    `json:"member_count,omitempty" db:"member_count"`
	MessageCount int    `json:"message_count,omitempty" db:"message_count"`
}

type ChatMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Joined fields
	User *UserPublic `json:"user,omitempty" db:"-"`
}

// ─── Notification ──────────────────────────────────────────────────────────

type Notification struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	UserID    uuid.UUID              `json:"user_id" db:"user_id"`
	Type      string                 `json:"type" db:"type"`
	Title     string                 `json:"title" db:"title"`
	Body      string                 `json:"body" db:"body"`
	IsRead    bool                   `json:"is_read" db:"is_read"`
	Metadata  map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

// ─── User Settings ─────────────────────────────────────────────────────────

type UserSettings struct {
	ID                  uuid.UUID              `json:"id" db:"id"`
	UserID              uuid.UUID              `json:"user_id" db:"user_id"`
	NotificationPrefs   map[string]interface{} `json:"notification_prefs" db:"notification_prefs"`
	Theme               string                 `json:"theme" db:"theme"`
	Language            string                 `json:"language" db:"language"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// ─── Request/Response DTOs ─────────────────────────────────────────────────

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=30,alphanum"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type AddLibraryRequest struct {
	MangaID uuid.UUID `json:"manga_id" validate:"required"`
	Status  string    `json:"status" validate:"required,oneof=reading completed plan_to_read on_hold dropped"`
}

type UpdateProgressRequest struct {
	CurrentChapter int     `json:"current_chapter" validate:"min=0"`
	CurrentVolume  int     `json:"current_volume" validate:"min=0"`
	Status         string  `json:"status" validate:"omitempty,oneof=reading completed plan_to_read on_hold dropped"`
	Rating         *float64 `json:"rating" validate:"omitempty,min=0,max=10"`
	Notes          string  `json:"notes" validate:"omitempty,max=2000"`
}

type SendMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ─── Stats ─────────────────────────────────────────────────────────────────

type UserStats struct {
	TotalManga      int     `json:"total_manga"`
	Reading         int     `json:"reading"`
	Completed       int     `json:"completed"`
	PlanToRead      int     `json:"plan_to_read"`
	OnHold          int     `json:"on_hold"`
	Dropped         int     `json:"dropped"`
	TotalChapters   int     `json:"total_chapters"`
	TotalVolumes    int     `json:"total_volumes"`
	AverageRating   float64 `json:"average_rating"`
	FavoriteGenres  []string `json:"favorite_genres"`
	MangaThisMonth  int     `json:"manga_this_month"`
	ChaptersThisMonth int   `json:"chapters_this_month"`
}

// ─── Server Health ─────────────────────────────────────────────────────────

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

// ─── WebSocket Events ──────────────────────────────────────────────────────

type WSEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserID  string      `json:"user_id,omitempty"`
	RoomID  string      `json:"room_id,omitempty"`
}
