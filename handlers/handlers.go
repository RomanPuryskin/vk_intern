package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vk_intern/internal/advertisements"
	"github.com/vk_intern/internal/logger"
	"github.com/vk_intern/internal/middleware"
	"github.com/vk_intern/internal/repository"
	"github.com/vk_intern/internal/users"
)

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param registerData body users.UserRequest true "Логин и пароль"
// @Success 201 {object} users.UserRegisterResponse
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /register [post]
func RegisterUser(c *fiber.Ctx) error {
	var newUser users.UserRequest

	//парсим JSON в структуру subscription
	if err := c.BodyParser(&newUser); err != nil {
		logger.L.Error("[RegisterUser | parse JSON]: failed parse newUser", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// валидация логина и пароля
	err := users.ValidateUserLoginPassword(&newUser)
	if err != nil {
		logger.L.Error("[RegisterUser | validate]:", "error", err)

		switch {
		case errors.Is(err, users.ErrShortLogin):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "логин должен содержать хотя бы 3 символа"})
		case errors.Is(err, users.ErrLongLogin):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "логин должен содержать не более 25 символов"})
		case errors.Is(err, users.ErrWrongLoginSymbols):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "логин может содержать только буквы различных алфавитов и цифры"})
		case errors.Is(err, users.ErrShortPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "пароль должен содержать хотя бы 8 символов"})
		case errors.Is(err, users.ErrLongPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "пароль должен содержать не более 25 символов"})
		case errors.Is(err, users.ErrWrongPasswordSymbols):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "пароль может содержать только буквы различных алфавитов и цифры"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// запрос к БД
	respUser, err := repository.RegisterUser(context.Background(), &newUser)
	if err != nil {
		switch {

		case errors.Is(err, repository.ErrUserExists):
			logger.L.Error("[RegisterUser | exec regiser]: login already exists", "login", newUser.Login)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "пользователь с таким логином уже существует"})

		default:
			logger.L.Error("[RegisterUser | exec regiser]:", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// успешный ответ
	logger.L.Info("[RegisterUser]: success RegisterUser request")
	return c.Status(fiber.StatusCreated).JSON(respUser)
}

// LoginUser godoc
// @Summary Авторизация пользователя
// @Description Авторизирует зарегистрированного пользователя
// @Security ApiKeyAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param loginData body users.UserRequest true "Логин и пароль"
// @Success 200 {string} string "token"
// @Success 208 {string} string "already authorized"
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /login [post]
func LoginUser(JWTsecret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var newUser users.UserRequest

		//парсим JSON в структуру subscription
		if err := c.BodyParser(&newUser); err != nil {
			logger.L.Error("[LoginUser | parse JSON]: failed parse newUser", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
		}

		// проверка введенных логина и пароля
		err := repository.CheckLoginAndPassword(context.Background(), &newUser)
		if err != nil {
			logger.L.Error("[LoginUser | validate]:", "error", err)

			switch {
			case errors.Is(err, repository.ErrUserLoginWrong):
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный логин пользователя"})
			case errors.Is(err, repository.ErrUserWrongPassword):
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "неверный пароль пользователя"})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		// создание JWT токена
		token, err := middleware.GenerateJWTToken(newUser.Login, JWTsecret)
		if err != nil {
			logger.L.Error("[LoginUser | generateJWT]:", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		logger.L.Info("[LoginUser]: success LoginUser request")
		return c.Status(fiber.StatusOK).JSON(token)
	}
}

// CreateAdvertisement godoc
// @Summary Создание объявления
// @Description Создает объявление, только для авторизированных пользователей
// @Security ApiKeyAuth
// @Tags advertisements
// @Accept json
// @Produce json
// @Param advertisementData body advertisements.CreateAdvertisementRequest true "Данные объявления"
// @Success 201 {object} advertisements.Advertisement
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 401 {object} map[string]interface{} "'error': 'unauthorized'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /advertisements [post]
func CreateAdvertisement(c *fiber.Ctx) error {
	var newAdv advertisements.CreateAdvertisementRequest

	//парсим JSON в структуру subscription
	if err := c.BodyParser(&newAdv); err != nil {
		logger.L.Error("[CreateAdvertisement | parse JSON]: failed parse newAdv", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// провалидируем данные
	err := advertisements.ValidateAdvertisement(&newAdv)
	if err != nil {
		logger.L.Error("[CreateAdvertisement | validate]:", "error", err)

		switch {
		case errors.Is(err, advertisements.ErrShortTitle):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "заголовок должен содержать хотя бы 3 символа"})
		case errors.Is(err, advertisements.ErrLongTitle):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "заголовк может содержать не более 50 символов"})
		case errors.Is(err, advertisements.ErrWrongTitleSymbols):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "заголовок может содержать только буквы и цифры"})
		case errors.Is(err, advertisements.ErrShortDescription):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "описание должно содержать хотя бы 1 символ"})
		case errors.Is(err, advertisements.ErrLongDescription):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "описание может содержать не более 500 символов"})
		case errors.Is(err, advertisements.ErrPriceLessZero):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может быть отрицательной"})
		case errors.Is(err, advertisements.ErrBigPrice):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может превышать 100 000 000"})
		case errors.Is(err, advertisements.ErrBigPricePrecision):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может содержать больше 2 знаков после запятой"})
		case errors.Is(err, advertisements.ErrWrongURL):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "некорректный URL"})
		case errors.Is(err, advertisements.ErrWrongImageFormat):
			formats := strings.Join(advertisements.ImageFormats, " ")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Поддерживаемые форматы: %s", formats)})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// получим логин из контекста
	loginInterface := c.Locals("login")
	if loginInterface == nil {
		logger.L.Error("[CreateAdvertisement | get login]: could not get login from token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	newAdv.UserLogin = loginInterface.(string)

	// запрос к БД
	respAdv, err := repository.LoadAdvertisement(context.Background(), &newAdv)
	if err != nil {
		logger.L.Error("[CreateAdvertisement | exec create adv]:", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.L.Info("[CreateAdvertisement]: success CreateAdvertisement request")
	return c.Status(fiber.StatusCreated).JSON(respAdv)
}

// GetAllAdvertisements godoc
// @Summary Показать список объявлений
// @Description Возвращает список объявлений по параметрам фильтрации
// @Security ApiKeyAuth
// @Tags advertisements
// @Accept json
// @Produce json
// @Param page query integer false "Номер страницы" default(1)
// @Param min_price query number false "Минимальная цена" default(0)
// @Param max_price query number false "Максимальная цена" default(100000000)
// @Param limit query integer false "Лимит на странице" default(10)
// @Param order_by query string false "Параметр для сортировки" default("created_at")
// @Param order query string false "Вид сортировки" default("DESC")
// @Success 200 {array} advertisements.AdvertisementResponse
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /advertisements [get]
func GetAllAdvertisements(c *fiber.Ctx) error {
	// получаем логин из контекста
	var login string
	loginInterface := c.Locals("login")
	if loginInterface == nil {
		login = ""
	} else {
		login = loginInterface.(string)
	}

	// получим параметры фильтрации из query
	params := advertisements.NewDefaultFilter()
	if err := c.QueryParser(&params); err != nil {
		logger.L.Error("[GetAllAdvertisements | parse JSON]: failed parse params", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "некорректные параметры фильтра"})
	}

	// провалидируем цены фильтрации
	err := advertisements.ValidatePricesInAdverisementFilter(&params)
	if err != nil {
		logger.L.Error("[GetAllAdvertisements | validate]:", "error", err)

		switch {
		case errors.Is(err, advertisements.ErrPriceLessZero):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может быть отрицательной"})
		case errors.Is(err, advertisements.ErrBigPrice):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может превышать 100 000 000"})
		case errors.Is(err, advertisements.ErrBigPricePrecision):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "цена не может содержать больше 2 знаков после запятой"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// запрос к БД
	advs, err := repository.GetAllAdvertisements(context.Background(), login, &params)
	if err != nil {
		logger.L.Error("[GetAllAdvertisements | exec get all advs]:", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.L.Info("[GetAllAdvertisements]: success GetAllAdvertisements request")
	return c.Status(fiber.StatusOK).JSON(advs)
}
