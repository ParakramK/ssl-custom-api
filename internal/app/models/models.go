package models

// All returns every app-owned model from a single source
func All() []any {
	return []any{
		&Module{},
		&User{},
		&ApiKey{},
		// &ModulePermission{},
	}
}
