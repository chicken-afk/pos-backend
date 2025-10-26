package response

import "pos/auth_service/app/entities"

type LoginResponse struct {
	SessionID    string    `json:"session_id"`
	User         UserLogin `json:"user"`
	TokenType    string    `json:"token_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    string    `json:"expires_in"`
}

type UserLogin struct {
	Email  string          `json:"email"`
	Name   string          `json:"name"`
	Role   entities.Role   `json:"role"`
	Outlet entities.Outlet `json:"outlet"`
}
