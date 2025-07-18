package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8082/lms/v1/auth"
)

var (
	client *resty.Client
)

// Request/Response structs
type RegisterRequest struct {
	Username              string `json:"username"`
	Email                 string `json:"email"`
	Password              string `json:"password"`
	Address               string `json:"address"`
	PhoneNumber           string `json:"phoneNumber"`
	MfaVerificationOption bool   `json:"mfaVerificationOption"`
}

type RegisterResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func init() {
	client = resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetHeader("Accept", "application/json")
	client.SetHeader("Content-Type", "application/json")

	// Enable debug mode for detailed request/response logging
	client.SetDebug(true)
}

func TestRegister(t *testing.T) {
	t.Run("Successful Registration", func(t *testing.T) {
		registerReq := RegisterRequest{
			Username:              "SwanhtetAungphyo",
			Email:                 "swanhtetaungp@gmail.com",
			Password:              "SwanHtet12@",
			Address:               "Swan Krakow",
			PhoneNumber:           "48608422691",
			MfaVerificationOption: true,
		}

		var registerResp RegisterResponse
		var errorResp ErrorResponse

		resp, err := client.R().
			SetBody(registerReq).
			SetResult(&registerResp).
			SetError(&errorResp).
			Post(baseURL + "/register")

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		assert.True(t, registerResp.Success)
		assert.NotEmpty(t, registerResp.ID)
		assert.NotEmpty(t, registerResp.Message)
	})

	t.Run("Registration with Invalid Email", func(t *testing.T) {
		registerReq := RegisterRequest{
			Username:              "TestUser",
			Email:                 "invalid-email",
			Password:              "Password123@",
			Address:               "Test Address",
			PhoneNumber:           "1234567890",
			MfaVerificationOption: false,
		}

		var errorResp ErrorResponse

		resp, err := client.R().
			SetBody(registerReq).
			SetError(&errorResp).
			Post(baseURL + "/register")

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.NotEmpty(t, errorResp.Error)
		assert.Contains(t, errorResp.Message, "email")
	})

	t.Run("Registration with Weak Password", func(t *testing.T) {
		registerReq := RegisterRequest{
			Username:              "TestUser",
			Email:                 "test@example.com",
			Password:              "123",
			Address:               "Test Address",
			PhoneNumber:           "1234567890",
			MfaVerificationOption: false,
		}

		var errorResp ErrorResponse

		resp, err := client.R().
			SetBody(registerReq).
			SetError(&errorResp).
			Post(baseURL + "/register")

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.NotEmpty(t, errorResp.Error)
		assert.Contains(t, errorResp.Message, "password")
	})

	t.Run("Registration with Missing Fields", func(t *testing.T) {
		registerReq := RegisterRequest{
			Username: "TestUser",
			// Missing email, password, etc.
		}

		var errorResp ErrorResponse

		resp, err := client.R().
			SetBody(registerReq).
			SetError(&errorResp).
			Post(baseURL + "/register")

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
		assert.NotEmpty(t, errorResp.Error)
	})

	t.Run("Duplicate Email Registration", func(t *testing.T) {
		registerReq := RegisterRequest{
			Username:              "AnotherUser",
			Email:                 "swanhtetaungp@gmail.com",
			Password:              "AnotherPassword123@",
			Address:               "Another Address",
			PhoneNumber:           "9876543210",
			MfaVerificationOption: false,
		}

		var errorResp ErrorResponse

		resp, err := client.R().
			SetBody(registerReq).
			SetError(&errorResp).
			Post(baseURL + "/register")

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode())
		assert.NotEmpty(t, errorResp.Error)
		assert.Contains(t, errorResp.Message, "email")
	})
}
