package users

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	minLenLogin    = 3
	maxLenLogin    = 25
	minLenPassword = 8
	maxLenPassword = 25
)

var (
	ErrShortLogin           = errors.New("login is shorter than required")
	ErrLongLogin            = errors.New("login is longer than required")
	ErrWrongLoginSymbols    = errors.New("wrong symbols in login")
	ErrShortPassword        = errors.New("password is shorter than required")
	ErrLongPassword         = errors.New("password is longer than required")
	ErrWrongPasswordSymbols = errors.New("wrong symbols in password")
)

func ValidateUserLoginPassword(user *UserRequest) error {
	if err := vaildateLogin(user.Login); err != nil {
		return fmt.Errorf("[ValidateUserLoginPassword] %w", err)
	}

	if err := validatePassword(user.Password); err != nil {
		return fmt.Errorf("[ValidateUserLoginPassword] %w", err)
	}

	return nil
}

func vaildateLogin(login string) error {
	// провалидируем логин на длину
	loginLen := utf8.RuneCountInString(login)
	if loginLen < minLenLogin {
		return fmt.Errorf("[validateLogin]: %w", ErrShortLogin)
	}
	if loginLen > maxLenLogin {
		return fmt.Errorf("[validateLogin]: %w", ErrLongLogin)
	}

	// провалидируем логин на содержание спец символов кроме букв и цифр
	if !regexp.MustCompile(`^[\p{L}\p{N}]+$`).MatchString(login) {
		return fmt.Errorf("[validateLogin]: %w", ErrWrongLoginSymbols)
	}
	return nil
}

func validatePassword(pass string) error {
	// провалидируем пароль на длину
	passLen := utf8.RuneCountInString(pass)
	if passLen < minLenPassword {
		return fmt.Errorf("[validatePassword]: %w", ErrShortPassword)
	}
	if passLen > maxLenPassword {
		return fmt.Errorf("[validatePassword]: %w", ErrLongPassword)
	}

	// провалидируем пароль на содержание спец символов кроме букв и цифр
	if !regexp.MustCompile(`^[\p{L}\p{N}]+$`).MatchString(pass) {
		return fmt.Errorf("[validatePassword]: %w", ErrWrongLoginSymbols)
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[HashPassword]: %w", err)
	}
	return string(hashed), nil
}

func ComparePasswordAndHashPassword(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
