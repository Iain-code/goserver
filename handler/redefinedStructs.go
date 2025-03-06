package handler

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Email          string         `json:"email"`
	HashedPassword sql.NullString `json:"-"`
}

type received struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt sql.NullTime `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

type CreateChirp struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

type incPost struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type RefreshTokenOnly struct {
	Token string `json:"token"`
}
