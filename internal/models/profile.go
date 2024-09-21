package models

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	ID        uuid.UUID `json:"id" db:"id"`                 // Unique id of Profile
	Email     string    `json:"email" db:"email"`           // User email
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Date and time when user created
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Date and time when user updated
}
