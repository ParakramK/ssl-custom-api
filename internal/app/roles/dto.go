package roles

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name    string `json:"name" doc:"Role name" minLength:"1" example:"admin"`
	IsAdmin bool   `json:"is_admin" doc:"Is admin role" example:"true"`
}
type CreateRoleInput struct {
	Body CreateRoleRequest
}

type CreateRoleResponse struct {
	ID      uuid.UUID `json:"id" doc:"Role ID"`
	Name    string    `json:"name" doc:"Role name"`
	IsAdmin bool      `json:"is_admin" doc:"Is admin role"`
}
type CreateRoleOutput struct {
	Body CreateRoleResponse
}

type RolesResponse struct {
	Roles []RoleRow `json:"roles" doc:"List of roles"`
}
type RolesOutput struct {
	Body []RolesResponse
}

type ListRolesRequest struct {
	Id   uuid.UUID `json:"id" doc:"Role ID"`
	Name string    `json:"name" doc:"Role name"`
}
