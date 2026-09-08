package utils

import (
	"github.com/google/uuid"
)

// NewV7ID returns a time-ordered UUIDv7, falling back to a random v4
// if the system RNG fails .

func NewV7ID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}
