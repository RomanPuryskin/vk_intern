package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vk_intern/internal/users"
)

var (
	ErrUserExists        = errors.New("user with this username already exists")
	ErrUserLoginWrong    = errors.New("user with this login does not exists")
	ErrUserWrongPassword = errors.New("wrong user password")
)

func RegisterUser(ctx context.Context, user *users.UserRequest) (*users.UserRegisterResponse, error) {
	// проверим существует ли уже пользователь с таким логином
	exists, err := checkLoginExists(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("[RegisterUser|check exists]: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("[RegisterUser|check exists]: %w", ErrUserExists)
	}

	//если не существует, добавляем
	hashedPassword, err := users.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("[RegisterUser|hash password] %w", err)
	}

	query := "INSERT INTO users (login,password,created_at) VALUES ($1 , $2 , $3)"
	created_at := time.Now()
	if _, err := Pool.Exec(ctx, query, user.Login, hashedPassword, created_at); err != nil {
		return nil, fmt.Errorf("[RegisterUser|exec register user]: %w", err)
	}

	// вернем данные добавленного пользователя
	var respUser users.UserRegisterResponse
	respUser.Login = user.Login
	respUser.Created_at = created_at

	return &respUser, nil
}

func CheckLoginAndPassword(ctx context.Context, user *users.UserRequest) error {
	// проверим существует ли пользователь с таким логином
	exists, err := checkLoginExists(ctx, user.Login)
	if err != nil {
		return fmt.Errorf("[CheckLoginAndPassword|check exists]: %w", err)
	}
	if !exists {
		return fmt.Errorf("[CheckLoginAndPassword|check exists]: %w", ErrUserLoginWrong)
	}

	// проверим корректность пароля
	if err := checkUserPasswordByLogin(ctx, user.Login, user.Password); err != nil {
		return fmt.Errorf("[CheckLoginAndPassword|check password]: %w", err)
	}
	return nil
}

func checkLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)"
	if err := Pool.QueryRow(ctx, query, login).Scan(&exists); err != nil {
		return false, fmt.Errorf("[checkLoginExists|exec check exists]: %w", err)
	}
	return exists, nil
}

func checkUserPasswordByLogin(ctx context.Context, login, password string) error {
	// получим пароль соответствующего пользователя
	var hashPass string
	query := "SELECT password FROM users WHERE login = $1"
	if err := Pool.QueryRow(ctx, query, login).Scan(&hashPass); err != nil {
		return fmt.Errorf("[checkUserPasswordByLogin|exec get password]: %w", err)
	}

	if !users.ComparePasswordAndHashPassword(hashPass, password) {
		return fmt.Errorf("[checkUserPasswordByLogin|compare passwords]: %w", ErrUserWrongPassword)
	}
	return nil
}
