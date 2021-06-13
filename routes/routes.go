package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kangana1024/goadmin-workshop/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/", controllers.Home)
	app.Post("/register", controllers.Register)
	app.Get("404", controllers.Handle404)
}
