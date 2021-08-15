package routes

import (
	"example.com/go/auth/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", controller.Hello)
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)
	app.Get("/api/user", controller.User)
	app.Post("/api/user", controller.Logout)

}
