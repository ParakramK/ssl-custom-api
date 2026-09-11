package users

import "github.com/google/uuid"

type UserRow struct {
	ID       uuid.UUID
	Username string
	RoleName string
}
