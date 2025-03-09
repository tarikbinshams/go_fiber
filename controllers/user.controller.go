package controllers

import (
	"context"
	"fiber/config"
	"fiber/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.DB.Collection("users")

	// Find all users
	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Println("Error decoding users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.JSON(users)
}
