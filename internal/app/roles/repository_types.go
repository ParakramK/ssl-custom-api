package roles

import "github.com/google/uuid"

type RoleRow struct {
	ID      uuid.UUID
	Name    string
	IsAdmin bool
}
