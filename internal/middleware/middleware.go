package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(JWTsecret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		tokenString := c.Get("Authorization")

		if tokenString != "" {
			err := checkTokenIsValid(tokenString, JWTsecret)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}

			return c.Status(fiber.StatusAlreadyReported).JSON("already authorized")
		}

		return c.Next()
	}
}

// строгий middleware, который не пускает дальше неавторизованных пользователей
func StrictMiddleware(JWTsecret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString != "" {
			err := checkTokenIsValid(tokenString, JWTsecret)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}

			//заносим логин из токена в контекст и переходим к роуту
			c.Locals("login", getLoginFromValidToken(tokenString, JWTsecret))
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
}

// нестрогий middleware, который пропускает и авторизованных и неавторизованных
func Middleware(JWTsecret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString != "" {
			err := checkTokenIsValid(tokenString, JWTsecret)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}

			//заносим логин из токена в контекст
			c.Locals("login", getLoginFromValidToken(tokenString, JWTsecret))

		}

		return c.Next()
	}
}
