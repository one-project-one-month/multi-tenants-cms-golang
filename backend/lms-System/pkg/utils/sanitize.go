package utils

import "strings"

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
