package main

import (
	"go-jwt/config"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New()

	log.Fatal(app.Listen(":3000"))
}
