package domain

type Session struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
}
