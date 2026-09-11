package modules

import "github.com/google/uuid"

type ModuleInfoResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}
type ModuleInfoOutput struct {
	Body ModuleInfoResponse `json:"body"`
}

type ModuleInfoRequest struct {
	ID   uuid.UUID `query:"id" doc:"Module ID"`
	Code string    `query:"code" doc:"Module code"`
}
type ModuleCreateRequest struct {
	Name        string  `json:"name" doc:"Module display name" minLength:"1" example:"Gatepass"`
	Code        string  `json:"code" doc:"Immutable module code" minLength:"1" example:"gatepass"`
	Description *string `json:"description,omitempty" doc:"Module description"`
}
type ModuleCreateInput struct {
	Body ModuleCreateRequest
}
type ResourceCreateRequest struct {
	Name        string    `json:"name" doc:"Resource display name" minLength:"1" example:"Sales"`
	Code        string    `json:"code" doc:"Immutable resource code" minLength:"1" example:"sales"`
	ModuleID    uuid.UUID `json:"module_id" doc:"Owning module ID"`
	Description *string   `json:"description,omitempty" doc:"Resource description"`
}
type ResourceCreateInput struct {
	Body ResourceCreateRequest
}

type ModuleCreateOutput struct {
	Body ModuleCreateResponse `json:"body"`
}

type ModuleCreateResponse struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResourceCreateOutput struct {
	Body ResourceCreateResponse `json:"body"`
}

type ResourceCreateResponse struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ListModulesRequest struct{}

type ListModulesOutput struct {
	Body ListModulesResponse `json:"body"`
}

type ListModulesResponse struct {
	Modules []ModuleInfoResponse `json:"modules"`
}
