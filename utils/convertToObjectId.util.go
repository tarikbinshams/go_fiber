package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectID(value string, fieldName string) (primitive.ObjectID, error) {
	if value == "" {
		return primitive.NilObjectID, fmt.Errorf("%s is required", fieldName)
	}
	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid %s format", fieldName)
	}
	return objectID, nil
}
