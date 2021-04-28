package controller

import "github.com/fedo3nik/cart-go-api/internal/domain/models"

type CreateCartResponse struct {
	ID    int               `json:"id"`
	Items []models.CartItem `json:"items"`
}
