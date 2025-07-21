package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/vk_intern/handlers"
	"github.com/vk_intern/internal/config"
	"github.com/vk_intern/internal/middleware"
)

func InitRoutes(app *fiber.App, cfg *config.Config) {
	auth := app.Group("/")
	auth.Post("/register", handlers.RegisterUser)
	auth.Post("/login", middleware.AuthMiddleware(cfg.JWT.JWTsecret), handlers.LoginUser(cfg.JWT.JWTsecret))

	adverts := app.Group("/advertisements")
	adverts.Post("/", middleware.StrictMiddleware(cfg.JWT.JWTsecret), handlers.CreateAdvertisement)
	adverts.Get("/", middleware.Middleware(cfg.JWT.JWTsecret), handlers.GetAllAdvertisements)

	app.Get("/swagger/*", swagger.HandlerDefault) // роут для сваггера
}
