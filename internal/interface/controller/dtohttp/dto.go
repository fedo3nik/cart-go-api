package controller

import "github.com/fedo3nik/cart-go-api/internal/domain/models"

type CreateCartResponse struct {
	ID    int               `json:"id"`
	Items []models.CartItem `json:"items"`
}

type AddItemRequest struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type AddItemResponse struct {
	ID       int    `json:"id"`
	CartID   int    `json:"cart_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
