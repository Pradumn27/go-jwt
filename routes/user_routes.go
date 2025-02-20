package routes

import (
	"go-jwt/config"
	"go-jwt/handlers"
	"go-jwt/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	auth := app.Group("/user")
	auth.Get("/details",middleware.JWTAuthGuard(config.DB), handlers.GetDetails)
}
