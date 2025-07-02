package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var jwtSecret = GetEnv("JWT_SECRET", "")

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
