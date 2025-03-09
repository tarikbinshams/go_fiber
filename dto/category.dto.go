package dto

type CategoryDTO struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}
