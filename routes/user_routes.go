package routes

import (
	"go-jwt/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	auth := app.Group("/user")
	auth.Get("/details", handlers.GetDetails)
}
