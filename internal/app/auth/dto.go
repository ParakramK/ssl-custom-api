package auth

type LoginRequest struct {
	Username string `json:"username" doc:"Username" minLength:"1" example:"admin"`
	Password string `json:"password" doc:"Password" minLength:"1" example:"secret"`
}

type LoginResponse struct {
	Token        string `json:"token" doc:"JWT access token"`
	RefreshToken string `json:"refresh_token" doc:"Refresh token"`
}

type LoginInput struct {
	Body LoginRequest
}

type LoginOutput struct {
	Body LoginResponse
}
