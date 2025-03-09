package routes

import (
	"fiber/controllers"
	"fiber/dto"
	"fiber/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", middlewares.ValidateBody[dto.UserRegisterDTO](), controllers.Register)
	auth.Post("/login", middlewares.ValidateBody[dto.UserLoginDTO](), controllers.Login)

	api.Use(middlewares.AuthMiddleware)
	api.Get("/users", controllers.GetUsers)

	category := api.Group("/categories")

	category.Post("/", middlewares.ValidateBody[dto.CategoryDTO](), controllers.CreateCategory)
	category.Get("/", controllers.GetCategories)
	category.Get("/:id", controllers.GetCategory)
	category.Patch("/:id", middlewares.ValidateBody[dto.CategoryDTO](), controllers.UpdateCategory)
	category.Delete("/:id", controllers.DeleteCategory)

	product := api.Group("/products")

	product.Post("/", middlewares.ValidateBody[dto.ProductDTO](), controllers.CreateProduct)
	product.Get("/", controllers.GetProducts)
	product.Get("/:id", controllers.GetProduct)
	product.Patch("/:id", middlewares.ValidateBody[dto.ProductDTO](), controllers.UpdateProduct)

}
