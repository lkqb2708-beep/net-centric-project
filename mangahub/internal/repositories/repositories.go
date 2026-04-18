package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"mangahub/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, bio, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`
	u.ID = uuid.New()
	return r.db.QueryRowContext(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.AvatarURL, u.Bio, u.Role,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id,username,email,password_hash,avatar_url,bio,role,is_active,last_seen_at,created_at,updated_at
		 FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.IsActive,
		&u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id,username,email,password_hash,avatar_url,bio,role,is_active,last_seen_at,created_at,updated_at
		 FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.IsActive,
		&u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id,username,email,password_hash,avatar_url,bio,role,is_active,last_seen_at,created_at,updated_at
		 FROM users WHERE username=$1`, username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.IsActive,
		&u.LastSeenAt, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET username=$1, avatar_url=$2, bio=$3, updated_at=NOW()
		 WHERE id=$4`,
		u.Username, u.AvatarURL, u.Bio, u.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`,
		hash, id,
	)
	return err
}

func (r *UserRepository) UpdateLastSeen(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET last_seen_at=$1 WHERE id=$2`, now, id,
	)
	return err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email,
	).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username,
	).Scan(&exists)
	return exists, err
}

func (r *UserRepository) CreateSettings(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_settings (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`,
		userID,
	)
	return err
}

// ─── Manga Repository ──────────────────────────────────────────────────────

type MangaRepository struct {
	db *sql.DB
}

func NewMangaRepository(db *sql.DB) *MangaRepository {
	return &MangaRepository{db: db}
}

func (r *MangaRepository) scanManga(row interface{ Scan(...interface{}) error }) (*models.Manga, error) {
	m := &models.Manga{}
	err := row.Scan(
		&m.ID, &m.Title, &m.Author, &m.Artist,
		pq.Array(&m.Genres), &m.Status,
		&m.ChapterCount, &m.VolumeCount,
		&m.Description, &m.CoverURL,
		&m.Year, &m.Rating, &m.PopularityRank,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MangaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Manga, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,title,author,artist,genres,status,chapter_count,volume_count,
		        description,cover_url,year,rating,popularity_rank,created_at,updated_at
		 FROM manga WHERE id=$1`, id)
	m, err := r.scanManga(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *MangaRepository) Search(ctx context.Context, query, status, genre string, page, pageSize int) ([]*models.Manga, int, error) {
	args := []interface{}{}
	where := "WHERE 1=1"
	i := 1

	if query != "" {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR author ILIKE $%d)", i, i+1)
		like := "%" + query + "%"
		args = append(args, like, like)
		i += 2
	}
	if status != "" {
		where += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, status)
		i++
	}
	if genre != "" {
		where += fmt.Sprintf(" AND $%d=ANY(genres)", i)
		args = append(args, genre)
		i++
	}

	// Count total
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM manga "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	orderBy := " ORDER BY popularity_rank ASC, rating DESC"
	limitClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id,title,author,artist,genres,status,chapter_count,volume_count,
		        description,cover_url,year,rating,popularity_rank,created_at,updated_at
		 FROM manga `+where+orderBy+limitClause, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var manga []*models.Manga
	for rows.Next() {
		m, err := r.scanManga(rows)
		if err != nil {
			return nil, 0, err
		}
		manga = append(manga, m)
	}
	return manga, total, rows.Err()
}

func (r *MangaRepository) GetPopular(ctx context.Context, limit int) ([]*models.Manga, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,title,author,artist,genres,status,chapter_count,volume_count,
		        description,cover_url,year,rating,popularity_rank,created_at,updated_at
		 FROM manga ORDER BY popularity_rank ASC, rating DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manga []*models.Manga
	for rows.Next() {
		m, err := r.scanManga(rows)
		if err != nil {
			return nil, err
		}
		manga = append(manga, m)
	}
	return manga, rows.Err()
}

func (r *MangaRepository) GetAll(ctx context.Context, page, pageSize int) ([]*models.Manga, int, error) {
	return r.Search(ctx, "", "", "", page, pageSize)
}

func (r *MangaRepository) Upsert(ctx context.Context, m *models.Manga) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO manga (id,title,author,artist,genres,status,chapter_count,volume_count,
		                    description,cover_url,year,rating,popularity_rank)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 ON CONFLICT (id) DO UPDATE SET
		   title=EXCLUDED.title, author=EXCLUDED.author, artist=EXCLUDED.artist,
		   genres=EXCLUDED.genres, status=EXCLUDED.status,
		   chapter_count=EXCLUDED.chapter_count, volume_count=EXCLUDED.volume_count,
		   description=EXCLUDED.description, cover_url=EXCLUDED.cover_url,
		   year=EXCLUDED.year, rating=EXCLUDED.rating,
		   popularity_rank=EXCLUDED.popularity_rank, updated_at=NOW()`,
		m.ID, m.Title, m.Author, m.Artist,
		pq.Array(m.Genres), m.Status,
		m.ChapterCount, m.VolumeCount,
		m.Description, m.CoverURL,
		m.Year, m.Rating, m.PopularityRank,
	)
	return err
}

// ─── Library Repository ────────────────────────────────────────────────────

type LibraryRepository struct {
	db *sql.DB
}

func NewLibraryRepository(db *sql.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) scanEntry(row interface{ Scan(...interface{}) error }) (*models.LibraryEntry, error) {
	e := &models.LibraryEntry{}
	err := row.Scan(
		&e.ID, &e.UserID, &e.MangaID, &e.Status,
		&e.CurrentChapter, &e.CurrentVolume,
		&e.Rating, &e.Notes,
		&e.StartedAt, &e.FinishedAt,
		&e.CreatedAt, &e.UpdatedAt,
	)
	return e, err
}

func (r *LibraryRepository) GetByUser(ctx context.Context, userID uuid.UUID, status string) ([]*models.LibraryEntry, error) {
	query := `
		SELECT le.id,le.user_id,le.manga_id,le.status,
		       le.current_chapter,le.current_volume,le.rating,le.notes,
		       le.started_at,le.finished_at,le.created_at,le.updated_at,
		       m.id,m.title,m.author,m.artist,m.genres,m.status,
		       m.chapter_count,m.volume_count,m.description,m.cover_url,
		       m.year,m.rating,m.popularity_rank,m.created_at,m.updated_at
		FROM library_entries le
		JOIN manga m ON m.id=le.manga_id
		WHERE le.user_id=$1`
	args := []interface{}{userID}
	if status != "" {
		query += " AND le.status=$2"
		args = append(args, status)
	}
	query += " ORDER BY le.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.LibraryEntry
	for rows.Next() {
		e := &models.LibraryEntry{Manga: &models.Manga{}}
		err := rows.Scan(
			&e.ID, &e.UserID, &e.MangaID, &e.Status,
			&e.CurrentChapter, &e.CurrentVolume,
			&e.Rating, &e.Notes,
			&e.StartedAt, &e.FinishedAt,
			&e.CreatedAt, &e.UpdatedAt,
			&e.Manga.ID, &e.Manga.Title, &e.Manga.Author, &e.Manga.Artist,
			pq.Array(&e.Manga.Genres), &e.Manga.Status,
			&e.Manga.ChapterCount, &e.Manga.VolumeCount,
			&e.Manga.Description, &e.Manga.CoverURL,
			&e.Manga.Year, &e.Manga.Rating, &e.Manga.PopularityRank,
			&e.Manga.CreatedAt, &e.Manga.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *LibraryRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.LibraryEntry, error) {
	e := &models.LibraryEntry{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id,user_id,manga_id,status,current_chapter,current_volume,
		        rating,notes,started_at,finished_at,created_at,updated_at
		 FROM library_entries WHERE id=$1 AND user_id=$2`, id, userID,
	).Scan(
		&e.ID, &e.UserID, &e.MangaID, &e.Status,
		&e.CurrentChapter, &e.CurrentVolume,
		&e.Rating, &e.Notes,
		&e.StartedAt, &e.FinishedAt,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *LibraryRepository) Add(ctx context.Context, e *models.LibraryEntry) error {
	e.ID = uuid.New()
	now := time.Now()
	if e.Status == "reading" {
		e.StartedAt = &now
	}
	return r.db.QueryRowContext(ctx,
		`INSERT INTO library_entries (id,user_id,manga_id,status,current_chapter,current_volume,notes,started_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING created_at,updated_at`,
		e.ID, e.UserID, e.MangaID, e.Status,
		e.CurrentChapter, e.CurrentVolume, e.Notes, e.StartedAt,
	).Scan(&e.CreatedAt, &e.UpdatedAt)
}

func (r *LibraryRepository) UpdateProgress(ctx context.Context, e *models.LibraryEntry) error {
	now := time.Now()
	if e.Status == "reading" && e.StartedAt == nil {
		e.StartedAt = &now
	}
	if e.Status == "completed" {
		e.FinishedAt = &now
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE library_entries
		 SET status=$1, current_chapter=$2, current_volume=$3,
		     rating=$4, notes=$5, started_at=COALESCE(started_at,$6),
		     finished_at=$7, updated_at=NOW()
		 WHERE id=$8 AND user_id=$9`,
		e.Status, e.CurrentChapter, e.CurrentVolume,
		e.Rating, e.Notes, e.StartedAt, e.FinishedAt,
		e.ID, e.UserID,
	)
	return err
}

func (r *LibraryRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM library_entries WHERE id=$1 AND user_id=$2`, id, userID,
	)
	return err
}

func (r *LibraryRepository) GetStats(ctx context.Context, userID uuid.UUID) (*models.UserStats, error) {
	s := &models.UserStats{}
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  COUNT(*) FILTER (WHERE 1=1)                  AS total,
		  COUNT(*) FILTER (WHERE status='reading')      AS reading,
		  COUNT(*) FILTER (WHERE status='completed')    AS completed,
		  COUNT(*) FILTER (WHERE status='plan_to_read') AS plan_to_read,
		  COUNT(*) FILTER (WHERE status='on_hold')      AS on_hold,
		  COUNT(*) FILTER (WHERE status='dropped')      AS dropped,
		  COALESCE(SUM(current_chapter),0)              AS total_chapters,
		  COALESCE(SUM(current_volume),0)               AS total_volumes,
		  COALESCE(AVG(rating) FILTER (WHERE rating IS NOT NULL),0) AS avg_rating
		FROM library_entries WHERE user_id=$1`, userID,
	).Scan(
		&s.TotalManga, &s.Reading, &s.Completed,
		&s.PlanToRead, &s.OnHold, &s.Dropped,
		&s.TotalChapters, &s.TotalVolumes, &s.AverageRating,
	)
	return s, err
}

// ─── History Repository ────────────────────────────────────────────────────

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) Add(ctx context.Context, h *models.ReadingHistory) error {
	h.ID = uuid.New()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO reading_history (id,user_id,manga_id,chapter_number,volume_number)
		 VALUES ($1,$2,$3,$4,$5) RETURNING read_at`,
		h.ID, h.UserID, h.MangaID, h.ChapterNumber, h.VolumeNumber,
	).Scan(&h.ReadAt)
}

func (r *HistoryRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.ReadingHistory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rh.id,rh.user_id,rh.manga_id,rh.chapter_number,rh.volume_number,rh.read_at,
		       m.id,m.title,m.author,m.cover_url
		FROM reading_history rh
		JOIN manga m ON m.id=rh.manga_id
		WHERE rh.user_id=$1
		ORDER BY rh.read_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.ReadingHistory
	for rows.Next() {
		h := &models.ReadingHistory{Manga: &models.Manga{}}
		err := rows.Scan(
			&h.ID, &h.UserID, &h.MangaID, &h.ChapterNumber, &h.VolumeNumber, &h.ReadAt,
			&h.Manga.ID, &h.Manga.Title, &h.Manga.Author, &h.Manga.CoverURL,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

// ─── Chat Repository ───────────────────────────────────────────────────────

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetRooms(ctx context.Context) ([]*models.ChatRoom, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cr.id,cr.name,cr.type,cr.manga_id,cr.description,cr.created_at,
		       (SELECT COUNT(*) FROM chat_messages WHERE room_id=cr.id) AS message_count
		FROM chat_rooms cr
		ORDER BY cr.type DESC, cr.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.ChatRoom
	for rows.Next() {
		room := &models.ChatRoom{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Type, &room.MangaID, &room.Description, &room.CreatedAt,
			&room.MessageCount,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *ChatRepository) GetMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*models.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cm.id,cm.room_id,cm.user_id,cm.content,cm.created_at,
		       u.id,u.username,u.avatar_url
		FROM chat_messages cm
		JOIN users u ON u.id=cm.user_id
		WHERE cm.room_id=$1
		ORDER BY cm.created_at ASC
		LIMIT $2 OFFSET $3`, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.ChatMessage
	for rows.Next() {
		m := &models.ChatMessage{User: &models.UserPublic{}}
		err := rows.Scan(
			&m.ID, &m.RoomID, &m.UserID, &m.Content, &m.CreatedAt,
			&m.User.ID, &m.User.Username, &m.User.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *models.ChatMessage) error {
	msg.ID = uuid.New()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO chat_messages (id,room_id,user_id,content)
		 VALUES ($1,$2,$3,$4) RETURNING created_at`,
		msg.ID, msg.RoomID, msg.UserID, msg.Content,
	).Scan(&msg.CreatedAt)
}

// ─── Notification Repository ────────────────────────────────────────────────

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*models.Notification, error) {
	query := `SELECT id,user_id,type,title,body,is_read,created_at
	          FROM notifications WHERE user_id=$1`
	if unreadOnly {
		query += " AND is_read=FALSE"
	}
	query += " ORDER BY created_at DESC LIMIT 50"

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, rows.Err()
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read=TRUE WHERE id=$1 AND user_id=$2`, id, userID,
	)
	return err
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read=TRUE WHERE user_id=$1 AND is_read=FALSE`, userID,
	)
	return err
}

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	n.ID = uuid.New()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO notifications (id,user_id,type,title,body)
		 VALUES ($1,$2,$3,$4,$5) RETURNING created_at`,
		n.ID, n.UserID, n.Type, n.Title, n.Body,
	).Scan(&n.CreatedAt)
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=FALSE`, userID,
	).Scan(&count)
	return count, err
}
