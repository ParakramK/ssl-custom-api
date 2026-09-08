package models

// All returns every app-owned model from a single source
func All() []any {
	return []any{
		&Role{},
		&User{},
		&Module{},
		&ApiKey{},
		// &APIPermission{},
		// &UserModulePermission{},
	}

}
