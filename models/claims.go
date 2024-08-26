package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// Claims struct to hold JWT claims with role information
type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}
