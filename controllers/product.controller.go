package controllers

import (
	"context"
	"fiber/config"
	"fiber/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productCollection := config.DB.Collection("products")
	categoryCollection := config.DB.Collection("categories")

	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Convert ObjectIDs
	// categoryObjectID, objectIdErr := utils.ConvertToObjectID(c.FormValue("category_id"), "category_id")
	// if objectIdErr != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": objectIdErr.Error()})
	// }

	// fmt.Println(categoryObjectID)

	var category models.Category
	catErr := categoryCollection.FindOne(ctx, bson.M{"_id": product.CategoryID}).Decode(&category)
	if catErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Category not found"})
	}

	// Get userID from context
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	// product.CategoryID = categoryObjectID
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	product.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	userObjectID, idErr := primitive.ObjectIDFromHex(userID)
	if idErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	product.CreatedBy = userObjectID
	product.UpdatedBy = userObjectID

	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product created successfully"})
}

func GetProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productCollection := config.DB.Collection("products")

	// MongoDB Aggregation Pipeline
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "category_id"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "created_by"},
			{"foreignField", "_id"},
			{"as", "created_by"},
		}}},
		{{"$unwind", bson.D{{"path", "$category"}, {"preserveNullAndEmptyArrays", true}}}},
		{{"$unwind", bson.D{{"path", "$created_by"}, {"preserveNullAndEmptyArrays", true}}}},
		// Exclude the password field from the user object
		{{"$project", bson.D{
			{"created_by.password", 0}, // Exclude password from created_by (user)
		}}},
	}

	// Execute Aggregation Query
	cursor, err := productCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	var products []bson.M // Use bson.M to handle dynamic structure
	if err := cursor.All(ctx, &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode products"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": products})
}

func GetProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productCollection := config.DB.Collection("products")

	// Get the ID from the URL parameters
	id := c.Params("id")

	// Convert the ID from string to ObjectID
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID format"})
	}

	// MongoDB Aggregation Pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", productID}}}},
		{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "category_id"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "created_by"},
			{"foreignField", "_id"},
			{"as", "created_by"},
		}}},
		{{"$unwind", bson.D{{"path", "$category"}, {"preserveNullAndEmptyArrays", true}}}},
		{{"$unwind", bson.D{{"path", "$created_by"}, {"preserveNullAndEmptyArrays", true}}}},
		// Exclude the password field from the user object
		{{"$project", bson.D{
			{"category.password", 0},   // Exclude category-related fields if needed
			{"created_by.password", 0}, // Exclude password from created_by (user)
		}}},
	}

	// Execute Aggregation Query
	cursor, err := productCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	var products []bson.M // Use bson.M to handle dynamic structure
	if err := cursor.All(ctx, &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode product"})
	}

	// If no product is found
	if len(products) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": products[0]})
}

func UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productCollection := config.DB.Collection("products")

	// Get the ID from the URL parameters
	id := c.Params("id")

	// Convert the ID from string to ObjectID
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID format"})
	}

	// Get userID from context
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	// Get the product from the request body
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userObjectID, idErr := primitive.ObjectIDFromHex(userID)
	if idErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Set the updated fields
	updateFields := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"image":       product.Image,
			"category_id": product.CategoryID,
			"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
			"updated_by":  userObjectID, // Use the logged-in user ID to track who updated
		},
	}

	// Update the product
	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": productID}, updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product updated successfully"})
}
