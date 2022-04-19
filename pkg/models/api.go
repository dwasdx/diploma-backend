package models

// JWTObject - structure of json web token response
type JWTObject struct {
	UserID    string `json:"user_id"`
	Access    string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}
