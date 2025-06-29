package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// IsEmail checks if input is a valid email address
func IsEmail(input string) bool {
	_, err := mail.ParseAddress(input)
	return err == nil
}

// IsMobile checks if input is a valid mobile number (basic check)
func IsMobile(input string) bool {
	// Remove spaces, dashes, plus signs
	cleaned := strings.ReplaceAll(input, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.TrimPrefix(cleaned, "+")

	// Mobile number: digits only, length between 7 and 15 (adjust as needed)
	matched, _ := regexp.MatchString(`^\d{7,15}$`, cleaned)
	return matched
}

// DetectLoginType returns "email" or "mobile" depending on input format
func DetectLoginType(login string) (string, error) {
	if IsEmail(login) {
		return "email", nil
	}
	if IsMobile(login) {
		return "mobile", nil
	}
	return "", fmt.Errorf("invalid login format")
}
