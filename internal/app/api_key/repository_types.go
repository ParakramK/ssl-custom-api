package apikey

import (
	"github.com/google/uuid"
)

type ApiKeyRow struct {
	ID       uuid.UUID
	Key      string
	UserName string
}
