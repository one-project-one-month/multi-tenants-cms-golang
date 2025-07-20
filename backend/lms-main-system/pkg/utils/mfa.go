package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"image/png"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// MFAConfig holds configuration for MFA token generation
type MFAConfig struct {
	Issuer      string
	AccountName string
	Algorithm   otp.Algorithm
	Digits      otp.Digits
	Period      uint
}

func NewMFAConfig(issuer string, accountName string) *MFAConfig {
	return &MFAConfig{
		Issuer:      issuer,
		AccountName: accountName,
		Algorithm:   otp.AlgorithmSHA256,
		Digits:      otp.DigitsSix,
		Period:      30,
	}
}

// DefaultMFAConfig returns a default MFA configuration
func DefaultMFAConfig() *MFAConfig {
	return &MFAConfig{
		Issuer:      "LMS",
		AccountName: "",
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.DigitsSix,
		Period:      30,
	}
}

// MFASecret represents an MFA secret key
type MFASecret struct {
	Secret          string
	QRCodeURL       string
	QRCodeBase64Str string
	Config          *MFAConfig
}

// GenerateMFASecret creates a new MFA secret
func GenerateMFASecret(config *MFAConfig, email string) (*MFASecret, error) {
	if config == nil {
		config = DefaultMFAConfig()
	}
	config.AccountName = email

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.Issuer,
		AccountName: config.AccountName,
		Period:      config.Period,
		Algorithm:   config.Algorithm,
		Digits:      config.Digits,
		Rand:        rand.Reader,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate MFA secret: %w", err)
	}

	buf := &bytes.Buffer{}
	image, err := key.Image(200, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to generate MFA secret: %w", err)
	}
	err = png.Encode(buf, image)
	if err != nil {
		logrus.Errorf("failed to generate MFA secret: %w", err)
	}
	return &MFASecret{
		Secret:          key.Secret(),
		QRCodeURL:       key.URL(),
		Config:          config,
		QRCodeBase64Str: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

// GenerateToken generates a TOTP token for the given secret and time
func GenerateToken(secret string, timestamp time.Time) (string, error) {
	token, err := totp.GenerateCodeCustom(secret, timestamp, totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// GenerateCurrentToken generates a token for the current time
func GenerateCurrentToken(secret string) (string, error) {
	return GenerateToken(secret, time.Now())
}

// VerifyToken verifies if a token is valid for the given secret
func VerifyToken(secret, token string) bool {
	return VerifyTokenWithSkew(secret, token, 1)
}

// VerifyTokenWithSkew verifies a token with time skew tolerance
// skew: number of time steps to check before and after current time
func VerifyTokenWithSkew(secret, token string, skew uint) bool {
	valid, err := totp.ValidateCustom(token, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      skew,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}

	return valid
}

// ValidateSecret checks if a secret is valid
func ValidateSecret(secret string) error {
	_, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return fmt.Errorf("invalid secret: %w", err)
	}

	return nil
}

// GetQRCodeURL generates a QR code URL for an existing secret
func GetQRCodeURL(secret, issuer, accountName string) (string, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, accountName, secret, issuer))
	if err != nil {
		return "", fmt.Errorf("failed to create QR code URL: %w", err)
	}

	return key.URL(), nil
}

type MFAManager struct {
	config *MFAConfig
}

// NewMFAManager creates a new MFA manager
func NewMFAManager(config *MFAConfig) *MFAManager {
	if config == nil {
		config = DefaultMFAConfig()
	}
	return &MFAManager{config: config}
}

// GenerateSecret generates a new MFA secret
func (m *MFAManager) GenerateSecret(email string) (*MFASecret, error) {
	return GenerateMFASecret(m.config, email)
}

// GenerateToken generates a token for the current time
func (m *MFAManager) GenerateToken(secret string) (string, error) {
	return GenerateCurrentToken(secret)
}

// VerifyToken verifies a token
func (m *MFAManager) VerifyToken(secret, token string) bool {
	return VerifyToken(secret, token)
}
