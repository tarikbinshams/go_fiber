package main

import (
	"fiber/config"
	"fiber/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	// docs are generated by Swag CLI, you have to import them.
	// replace with your own docs folder, usually "github.com/username/reponame/docs"
	_ "fiber/docs"
)

func main() {
	app := fiber.New()

	config.ConnectDB()

	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)

	app.Listen(":4000")
}
