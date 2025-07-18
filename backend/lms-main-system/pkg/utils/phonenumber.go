package utils

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumberResult struct {
	E164Format          string // +48608422691
	InternationalFormat string // +48 608 422 691
	NationalFormat      string // 608 422 691
	CountryCode         int32  // 48
	NationalNumber      uint64 // 608422691
	IsValid             bool
	IsPossible          bool
	Region              string // PL
}

func ParsePhoneNumberEnhanced(phoneNumber string, defaultCountryCode string) (*PhoneNumberResult, error) {
	if len(defaultCountryCode) != 2 {
		return nil, fmt.Errorf("invalid country code: must be 2 characters")
	}

	parsedNumber, err := phonenumbers.Parse(phoneNumber, defaultCountryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse phone number: %w", err)
	}

	isValid := phonenumbers.IsValidNumber(parsedNumber)
	isPossible := phonenumbers.IsPossibleNumber(parsedNumber)

	if !isValid && !isPossible {
		return nil, fmt.Errorf("invalid phone number: not valid or possible")
	}

	region := phonenumbers.GetRegionCodeForNumber(parsedNumber)

	result := &PhoneNumberResult{
		E164Format:          phonenumbers.Format(parsedNumber, phonenumbers.E164),
		InternationalFormat: phonenumbers.Format(parsedNumber, phonenumbers.INTERNATIONAL),
		NationalFormat:      phonenumbers.Format(parsedNumber, phonenumbers.NATIONAL),
		CountryCode:         parsedNumber.GetCountryCode(),
		NationalNumber:      parsedNumber.GetNationalNumber(),
		IsValid:             isValid,
		IsPossible:          isPossible,
		Region:              region,
	}

	return result, nil
}
