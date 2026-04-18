package handlers

import (
	"net/http"
	"time"

	"mangahub/internal/auth"
	"mangahub/internal/middleware"
	"mangahub/internal/models"
	"mangahub/internal/repositories"
)

type AuthHandler struct {
	users *repositories.UserRepository
}

func NewAuthHandler(users *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

// POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "username, email, and password are required")
		return
	}
	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	ctx := r.Context()

	emailExists, _ := h.users.EmailExists(ctx, req.Email)
	if emailExists {
		respondError(w, http.StatusConflict, "email already registered")
		return
	}
	usernameExists, _ := h.users.UsernameExists(ctx, req.Username)
	if usernameExists {
		respondError(w, http.StatusConflict, "username already taken")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "user",
		AvatarURL:    "",
		Bio:          "",
	}
	if err := h.users.Create(ctx, user); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	_ = h.users.CreateSettings(ctx, user.ID)

	token, err := auth.GenerateToken(user.ID.String(), user.Username, user.Email, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	setTokenCookie(w, token)
	respondCreated(w, models.LoginResponse{
		Token: token,
		User:  toPublicUser(user),
	})
}

// POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()

	user, err := h.users.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !user.IsActive {
		respondError(w, http.StatusForbidden, "account is deactivated")
		return
	}

	token, err := auth.GenerateToken(user.ID.String(), user.Username, user.Email, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	_ = h.users.UpdateLastSeen(ctx, user.ID)
	setTokenCookie(w, token)
	respondOK(w, models.LoginResponse{
		Token: token,
		User:  toPublicUser(user),
	})
}

// POST /api/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
	respondOK(w, map[string]string{"message": "logged out"})
}

// GET /api/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	ctx := r.Context()
	user, err := h.users.GetByEmail(ctx, claims.Email)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	respondOK(w, toPublicUser(user))
}

// PUT /api/me
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	var req models.UpdateProfileRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	user, _ := h.users.GetByEmail(ctx, claims.Email)
	if user == nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if err := h.users.UpdateProfile(ctx, user); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	respondOK(w, toPublicUser(user))
}

// POST /api/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	var req models.ChangePasswordRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	user, _ := h.users.GetByEmail(ctx, claims.Email)
	if user == nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	if !auth.CheckPassword(req.OldPassword, user.PasswordHash) {
		respondError(w, http.StatusUnauthorized, "incorrect current password")
		return
	}
	if len(req.NewPassword) < 8 {
		respondError(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}
	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	if err := h.users.UpdatePassword(ctx, user.ID, hash); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update password")
		return
	}
	respondOK(w, map[string]string{"message": "password updated"})
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
	})
}

func toPublicUser(u *models.User) models.UserPublic {
	return models.UserPublic{
		ID:         u.ID,
		Username:   u.Username,
		AvatarURL:  u.AvatarURL,
		Bio:        u.Bio,
		Role:       u.Role,
		LastSeenAt: u.LastSeenAt,
		CreatedAt:  u.CreatedAt,
	}
}
