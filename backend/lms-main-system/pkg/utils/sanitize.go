package utils

import (
	"strings"
	"unicode"
)

func SanitizeUsername(username string) string {
	return strings.ToLower(username)
}

func SanitizeEmail(email string) string {
	return strings.ToLower(email)
}

func SanitizePhoneNumber(phoneNumber string) string {
	if len(phoneNumber) == 0 {
		return ""
	}
	if phoneNumber[0] == '+' {
		return strings.ReplaceAll(phoneNumber, "+", "")
	}
	return phoneNumber
}

func SanitizePhoneNumberEnhanced(phoneNumber string) string {
	if len(phoneNumber) == 0 {
		return ""
	}

	phoneNumber = strings.TrimSpace(phoneNumber)
	var result strings.Builder

	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		result.WriteRune('+')
	}

	for _, char := range phoneNumber {
		if unicode.IsDigit(char) {
			result.WriteRune(char)
		}
	}

	return result.String()
}
