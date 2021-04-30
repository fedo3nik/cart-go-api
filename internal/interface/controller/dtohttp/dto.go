package controller

import "github.com/fedo3nik/cart-go-api/internal/domain/models"

// A CartResponse represents json response for the CreateCart handler.
type CartResponse struct {
	ID    int               `json:"id"`    // Cart ID
	Items []models.CartItem `json:"items"` // Items in the cart
}

// An AddItemRequest represents json request for the AddItem handler.
type AddItemRequest struct {
	Product  string `json:"product"`  // Product title
	Quantity int    `json:"quantity"` // Quantity of the products in the item
}

// An AddItemResponse represents json request for the AddItem handler.
type AddItemResponse struct {
	ID       int    `json:"id"`       // Item ID
	CartID   int    `json:"cart_id"`  // ID of the cart in which item was placed
	Product  string `json:"product"`  // Product title
	Quantity int    `json:"quantity"` // Quantity of the products in the item
}

// A RemoveItemResponse represents json response for the RemoveItem handler.
type RemoveItemResponse struct {
}

// An ErrorResponse represents json response for the cases when the error is occurred.
type ErrorResponse struct {
	Message string `json:"message"` // Message of the error
}
