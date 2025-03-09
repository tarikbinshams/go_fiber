package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Image       string             `bson:"image" json:"image"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}
