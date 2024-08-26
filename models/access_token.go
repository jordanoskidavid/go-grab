package models

import "time"

type AccessToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
}
