// pkg/utils/error_response.go
package utils

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
	ErrorTypeAuthorization  ErrorType = "AUTHORIZATION_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict       ErrorType = "CONFLICT_ERROR"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
	ErrorTypeRateLimit      ErrorType = "RATE_LIMIT_ERROR"
	ErrorTypeExternal       ErrorType = "EXTERNAL_SERVICE_ERROR"
)

type ErrorResponse struct {
	Type      ErrorType `json:"type"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ToGRPCStatus converts ErrorResponse to gRPC status
func (e *ErrorResponse) ToGRPCStatus() error {
	var grpcCode codes.Code

	switch e.Type {
	case ErrorTypeValidation:
		grpcCode = codes.InvalidArgument
	case ErrorTypeAuthentication:
		grpcCode = codes.Unauthenticated
	case ErrorTypeAuthorization:
		grpcCode = codes.PermissionDenied
	case ErrorTypeNotFound:
		grpcCode = codes.NotFound
	case ErrorTypeConflict:
		grpcCode = codes.AlreadyExists
	case ErrorTypeRateLimit:
		grpcCode = codes.ResourceExhausted
	case ErrorTypeExternal:
		grpcCode = codes.Unavailable
	case ErrorTypeInternal:
		fallthrough
	default:
		grpcCode = codes.Internal
	}

	message := e.Message
	if e.Details != "" {
		message = fmt.Sprintf("%s: %s", e.Message, e.Details)
	}

	return status.Error(grpcCode, message)
}

func ErrInvalidCredentials(details ...string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthentication,
		Code:    "AUTH_INVALID_CREDENTIALS",
		Message: "Invalid email or password",
		Details: joinDetails(details...),
	}
}

func ErrUserNotFound(email string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeNotFound,
		Code:    "AUTH_USER_NOT_FOUND",
		Message: "User not found",
		Details: fmt.Sprintf("No user found with email: %s", email),
	}
}

func ErrUserAlreadyExists(email string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeConflict,
		Code:    "AUTH_USER_EXISTS",
		Message: "User already exists",
		Details: fmt.Sprintf("User with email %s already exists", email),
	}
}

func ErrInvalidToken(details ...string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthentication,
		Code:    "AUTH_INVALID_TOKEN",
		Message: "Invalid or expired token",
		Details: joinDetails(details...),
	}
}

func ErrTokenExpired() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthentication,
		Code:    "AUTH_TOKEN_EXPIRED",
		Message: "Token has expired",
	}
}

func ErrUnauthorized(details ...string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthorization,
		Code:    "AUTH_UNAUTHORIZED",
		Message: "Unauthorized access",
		Details: joinDetails(details...),
	}
}

func ErrInvalidEmail() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION_INVALID_EMAIL",
		Message: "Invalid email format",
	}
}

func ErrWeakPassword() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION_WEAK_PASSWORD",
		Message: "Password does not meet security requirements",
		Details: "Password must be at least 8 characters long and contain uppercase, lowercase, number, and special character",
	}
}

func ErrMissingField(field string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION_MISSING_FIELD",
		Message: "Required field is missing",
		Details: fmt.Sprintf("Field '%s' is required", field),
	}
}

func ErrInvalidFieldValue(field, value string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "VALIDATION_INVALID_FIELD",
		Message: "Invalid field value",
		Details: fmt.Sprintf("Invalid value '%s' for field '%s'", value, field),
	}
}

func ErrMissingOrganization() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "ORG_MISSING_HEADER",
		Message: "Organization context required",
		Details: "x-organisation header is required",
	}
}

func ErrOrganizationNotFound(orgName string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeNotFound,
		Code:    "ORG_NOT_FOUND",
		Message: "Organization not found",
		Details: fmt.Sprintf("Organization '%s' does not exist", orgName),
	}
}

// Rate Limiting Errors
func ErrRateLimitExceeded(action string, resetTime string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeRateLimit,
		Code:    "RATE_LIMIT_EXCEEDED",
		Message: fmt.Sprintf("Rate limit exceeded for %s", action),
		Details: fmt.Sprintf("Please try again after %s", resetTime),
	}
}

// MFA Errors
func ErrMFARequired() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthentication,
		Code:    "AUTH_MFA_REQUIRED",
		Message: "Multi-factor authentication required",
	}
}

func ErrInvalidMFACode() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeAuthentication,
		Code:    "AUTH_INVALID_MFA_CODE",
		Message: "Invalid MFA verification code",
	}
}

func ErrMFANotSetup() *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeValidation,
		Code:    "AUTH_MFA_NOT_SETUP",
		Message: "MFA is not set up for this user",
	}
}

// Internal Errors
func ErrInternal(details ...string) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeInternal,
		Code:    "INTERNAL_ERROR",
		Message: "An internal error occurred",
		Details: joinDetails(details...),
	}
}

func ErrDatabaseOperation(operation string, err error) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeInternal,
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
		Details: fmt.Sprintf("Failed to %s: %v", operation, err),
	}
}

func ErrEmailDelivery(email string, err error) *ErrorResponse {
	return &ErrorResponse{
		Type:    ErrorTypeExternal,
		Code:    "EMAIL_DELIVERY_ERROR",
		Message: "Failed to send email",
		Details: fmt.Sprintf("Could not send email to %s: %v", email, err),
	}
}

func joinDetails(details ...string) string {
	if len(details) == 0 {
		return ""
	}
	if len(details) == 1 {
		return details[0]
	}
	result := details[0]
	for _, detail := range details[1:] {
		result += "; " + detail
	}
	return result
}

func NewError(errorType ErrorType, code, message string, details ...string) *ErrorResponse {
	return &ErrorResponse{
		Type:    errorType,
		Code:    code,
		Message: message,
		Details: joinDetails(details...),
	}
}
