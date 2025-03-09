package dto

type ProductDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	// Image       string  `json:"image" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"` // Ensure JSON key is lowercase
}
