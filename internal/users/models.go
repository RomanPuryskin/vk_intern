package users

import "time"

// UserRequest модель запроса на регистрацию
// @Description Модель описывает запрос на регистрацию с логином и паролем
type UserRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserRegisterResponse модель ответа на регистрацию
// @Description Модель описывает ответ на успешную регистрацию
type UserRegisterResponse struct {
	Login      string    `json:"login"`
	Created_at time.Time `json:"created_at" example:"2023-05-15T10:00:00Z" format:"date-time"`
}
