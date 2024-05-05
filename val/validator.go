package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isUsernameValid = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isFullNameValid = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("string length must be between %d and %d", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 20); err != nil {
		return err
	}

	if !isUsernameValid(username) {
		return fmt.Errorf("username can only contain letters, numbers and underscores")
	}

	return nil
}

func ValidateFullName(username string) error {
	if err := ValidateString(username, 3, 20); err != nil {
		return err
	}

	if !isFullNameValid(username) {
		return fmt.Errorf("username can only contain letters and space")
	}

	return nil
}

func ValidatePassword(password string) error {
	if err := ValidateString(password, 6, 20); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 6, 50); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

func ValidateEmailId(emailId int64) error {
	if emailId <= 0 {
		return fmt.Errorf("email_id must be greater than 0")
	}

	return nil
}

func ValidateVerifyToken(token string) error {
	if err := ValidateString(token, 32, 32); err != nil {
		return err
	}

	return nil
}
