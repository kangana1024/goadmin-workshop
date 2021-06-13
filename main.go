package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kangana1024/goadmin-workshop/database"
	"github.com/kangana1024/goadmin-workshop/routes"
)
type User struct {
	Name string
}
func main() {

	err :=	database.Connect()
	if err != nil {
		panic("Could not connect Database!.")
	}

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":3000")
}