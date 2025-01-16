package entity

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID
	URL       string
	Alias     string
	ExpiredAt time.Time
}
