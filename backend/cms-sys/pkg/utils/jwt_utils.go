package utils

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"os"
	"time"
)

var jwtSecret = GetEnv("JWT_SECRET", "")
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

func InitJWTKeysFromVault() error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	if vaultAddr == "" || vaultToken == "" {
		return errors.New("VAULT_ADDR or VAULT_TOKEN not set")
	}

	client, err := vault.NewClient(&vault.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return err
	}
	client.SetToken(vaultToken)

	secret, err := client.KVv2("secret").Get(context.Background(), "jwt")
	if err != nil {
		return err
	}

	privatePEM := []byte(secret.Data["private_key"].(string))
	publicPEM := []byte(secret.Data["public_key"].(string))

	// Parse private key
	block, _ := pem.Decode(privatePEM)
	if block == nil {
		return errors.New("invalid private key PEM")
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	block, _ = pem.Decode(publicPEM)
	if block == nil {
		return errors.New("invalid public key PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var ok bool
	publicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	return nil
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
	return token.SignedString(privateKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
