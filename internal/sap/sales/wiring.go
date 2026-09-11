package sales

import (
	"ssl-custom-api/internal/providers/hana"
)

func New(provider *hana.Provider) *Handler {
	repo := NewRepository(provider)
	service := NewService(repo)

	return NewHandler(service)
}
