package users

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username string    `json:"username" doc:"Username" minLength:"5" example:"admin"`
	Password string    `json:"password" doc:"Password" minLength:"8" example:"secret"`
	Email    string    `json:"email" doc:"Email" minLength:"1" example:"admin@example.com"`
	RoleId   uuid.UUID `json:"role_id" doc:"Role ID" example:"00000000-0000-0000-0000-000000000000"`
}
type CreateUserInput struct {
	Body CreateUserRequest
}

type CreateUserResponse struct {
	ID       uuid.UUID `json:"id" doc:"User ID"`
	Username string    `json:"username" doc:"Username"`
}
type CreateUserOutput struct {
	Body CreateUserResponse
}
