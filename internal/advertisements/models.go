package advertisements

import "time"

// Advertisement полная модель объявления
// @Description Модель описывает объявление для возврата при его создании
type Advertisement struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	UserLogin   string    `json:"userlogin"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-15T10:00:00Z" format:"date-time"`
}

// Advertisement модель запроса объявления
// @Description Модель описывает запрос на создание объявления
type CreateAdvertisementRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	ImageURL    string  `json:"image_url" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	UserLogin   string
}

// Advertisement модель объявления при получении
// @Description Модель описывает ответ на получение объявления
type AdvertisementResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	UserLogin   string    `json:"userlogin"`
	CreatedAt   time.Time `json:"created_at" example:"2023-05-15T10:00:00Z" format:"date-time"`
	IsMine      bool      `json:"ismine,omitempty"`
}

type AdvertisementFilter struct {
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
	OrderBy  string  `query:"order_by" validate:"oneof=price created_at"`
	Order    string  `query:"order" validate:"oneof=ASC DESC"`
	MinPrice float64 `query:"min_price"`
	MaxPrice float64 `query:"max_price"`
}
