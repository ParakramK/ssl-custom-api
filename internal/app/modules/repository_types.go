package modules

import "github.com/google/uuid"

type ModuleRow struct {
	Id          uuid.UUID
	Name        string
	Code        string
	Description *string
}

type ResourceRow struct {
	Id          uuid.UUID
	Name        string
	Code        string
	Description *string
	ModuleName  string
}
