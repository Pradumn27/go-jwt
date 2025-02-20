package main

import (
	"go-jwt/config"
	"go-jwt/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New()

	routes.SetupAuthRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
