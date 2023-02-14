package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	ID uuid.UUID `json:"id"`
}
