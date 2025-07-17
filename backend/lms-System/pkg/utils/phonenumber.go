package utils

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParsePhoneNumber(phoneNumber string, countryCode string) (string, bool, error) {
	if len(countryCode) != 2 {
		return "", false, fmt.Errorf("invalid country code")
	}

	normalizedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse phone number: %s", err)
	}

	if !phonenumbers.IsValidNumber(normalizedPhoneNumber) &&
		!phonenumbers.IsPossibleNumber(normalizedPhoneNumber) &&
		!phonenumbers.IsValidNumberForRegion(normalizedPhoneNumber, countryCode) {
		return "", false, status.Errorf(codes.InvalidArgument, "invalid phone number")
	}

	return fmt.Sprintf("%d%d", normalizedPhoneNumber.GetCountryCode(), normalizedPhoneNumber.GetNationalNumber()), phonenumbers.IsValidNumber(normalizedPhoneNumber), nil
}
