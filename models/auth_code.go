package models

import "time"

type AuthorizationCode struct {
	Code      string
	UserID    int
	ExpiresAt time.Time
}
