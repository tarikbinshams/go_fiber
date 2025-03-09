package controllers

import (
	"context"
	"fiber/config"
	"fiber/models"
	"fiber/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func createToken(email string, userId primitive.ObjectID) (string, error) {
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, config.AuthClaims{
		Email:  email,
		UserId: userId.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		fmt.Printf("Error generating token string: %v", err)
		return "", err
	}
	return tokenString, nil
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.DB.Collection("users")

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := createToken(user.Email, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func Register(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Allowed file types & max size (5MB)
	acceptedTypes := []string{"jpg", "jpeg", "png"}
	maxSize := int64(5 * 1024 * 1024) // 5MB

	var image *multipart.FileHeader

	// Check if a file was uploaded
	if file, err := c.FormFile("image"); err == nil {
		image = file
		// Validate file only if it exists
		err := utils.ValidateFile(file, acceptedTypes, maxSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Store file (or move to a folder)
		log.Printf("Uploaded file: %s", file.Filename)
	}

	usersCollection := config.DB.Collection("users")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	user.ID = primitive.NewObjectID()
	user.Status = "active"

	// If file exists, store the filename
	if image != nil {
		user.Image = image.Filename
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password" + err.Error()})
	}
	user.Password = string(hashedPassword)

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(user)
}
