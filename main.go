package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kangana1024/goadmin-workshop/database"
	"github.com/kangana1024/goadmin-workshop/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	err = database.Connect()
	if err != nil {
		panic("Could not connect Database!.")
	}

	app := fiber.New()

	corsOrigin := cors.New(cors.Config{
		AllowCredentials: true,
	})

	app.Use(corsOrigin)

	routes.Setup(app)

	app.Listen(":3000")
}
