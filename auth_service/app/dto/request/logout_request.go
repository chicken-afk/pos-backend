package request

type LogoutRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}
