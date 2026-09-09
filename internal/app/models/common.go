package models

type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
	PermissionAdmin  Permission = "admin"
)

// AllPermissions lists every valid permission value.
func AllPermissions() []Permission {
	return []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionDelete,
		PermissionAdmin,
	}
}

// ValidPermission reports whether p is a known permission.
func ValidPermission(p Permission) bool {
	for _, known := range AllPermissions() {
		if p == known {
			return true
		}
	}
	return false
}
