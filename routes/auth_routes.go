package routes

import (
	"go-jwt/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/signup", handlers.SignUp)
	auth.Post("/signin", handlers.SignIn)
	auth.Get("/refresh",  handlers.RefreshToken)
	auth.Post("/logout", handlers.Logout)
}
