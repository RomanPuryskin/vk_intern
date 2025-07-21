package advertisements

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	minLenTitle       = 3
	maxLenTitle       = 50
	minLenDescription = 1
	maxLenDescription = 500
	maxPrice          = 100000000
	maxPricePrecision = 2
)

var ImageFormats = []string{".jpg", ".jpeg", ".png", ".webp"}

var (
	ErrShortTitle        = errors.New("title is shorter than required")
	ErrLongTitle         = errors.New("title is longer than required")
	ErrWrongTitleSymbols = errors.New("wrong symbols in title")
	ErrShortDescription  = errors.New("description is shorter than required")
	ErrLongDescription   = errors.New("description is longer than required")
	ErrPriceLessZero     = errors.New("price can not be less zero")
	ErrBigPrice          = errors.New("price bigger than required")
	ErrBigPricePrecision = errors.New("price precision bigger than required")
	ErrWrongImageFormat  = errors.New("wrong image format")
	ErrWrongURL          = errors.New("wrong image URL")
)

func ValidateAdvertisement(adv *CreateAdvertisementRequest) error {
	if err := validateTitle(adv.Title); err != nil {
		return fmt.Errorf("[ValidateAdvertisement|title] %w", err)
	}

	if err := vaildateDescription(adv.Description); err != nil {
		return fmt.Errorf("[ValidateAdvertisement|description] %w", err)
	}

	if err := validatePrice(adv.Price); err != nil {
		return fmt.Errorf("[ValidateAdvertisement|price] %w", err)
	}

	if err := validateImageURL(adv.ImageURL); err != nil {
		return fmt.Errorf("[ValidateAdvertisement|url] %w", err)
	}

	return nil
}

func validateTitle(title string) error {
	// провалидируем заголовок на длину
	titleLen := utf8.RuneCountInString(title)
	if titleLen < minLenTitle {
		return fmt.Errorf("[validateTitle]: %w", ErrShortTitle)
	}
	if titleLen > maxLenTitle {
		return fmt.Errorf("[validateTitle]: %w", ErrLongTitle)
	}

	// провалидируем заголовок на содержание спец символов кроме букв и цифр
	if !regexp.MustCompile(`^[\p{L}\p{N}]+$`).MatchString(title) {
		return fmt.Errorf("[validateTitle]: %w", ErrWrongTitleSymbols)
	}
	return nil
}

func vaildateDescription(desc string) error {
	// провалидируем описание на длину
	descLen := utf8.RuneCountInString(desc)
	if descLen < minLenDescription {
		return fmt.Errorf("[vaildateDescription]: %w", ErrShortDescription)
	}
	if descLen > maxLenDescription {
		return fmt.Errorf("[vaildateDescription]: %w", ErrLongDescription)
	}

	return nil
}

func validatePrice(price float64) error {
	if price < 0 {
		return fmt.Errorf("[validatePrice]: %w", ErrPriceLessZero)
	}
	if price > maxPrice {
		return fmt.Errorf("[validatePrice]: %w", ErrBigPrice)
	}

	// цену можно указать не более чем с двумя знаками после запятой(копейки)
	priceStr := fmt.Sprintf("%v", price)
	parts := strings.Split(priceStr, ".")
	if len(parts) == 2 && len(parts[1]) > maxPricePrecision {
		return fmt.Errorf("[validatePrice]: %w", ErrBigPricePrecision)
	}
	return nil
}

func validateImageURL(imURL string) error {
	u, err := url.Parse(imURL)
	if err != nil {
		return fmt.Errorf("[validateImageURL]: %w", ErrWrongURL)
	}

	// проверим формат
	ext := strings.ToLower(path.Ext(u.Path))
	for _, e := range ImageFormats {
		if e == ext {
			return nil
		}
	}
	return fmt.Errorf("[validateImageURL]: %w", ErrWrongImageFormat)
}

func NewDefaultFilter() AdvertisementFilter {
	return AdvertisementFilter{
		Page:     1,
		Limit:    10,
		OrderBy:  "created_at",
		Order:    "DESC",
		MinPrice: 0,
		MaxPrice: 100000000,
	}
}

func ValidatePricesInAdverisementFilter(filer *AdvertisementFilter) error {
	if err := validatePrice(filer.MinPrice); err != nil {
		return fmt.Errorf("[ValidatePricesInAdverisementFilter|min_price] %w", err)
	}

	if err := validatePrice(filer.MaxPrice); err != nil {
		return fmt.Errorf("[ValidatePricesInAdverisementFilter|max_price] %w", err)
	}
	return nil
}
