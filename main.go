package main

import (
	"example.com/go/auth/database"
	"example.com/go/auth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	//connect to database
	database.Connect()

	app := fiber.New()

	// setup routes

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":80")
}
