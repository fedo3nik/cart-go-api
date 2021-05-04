// Package doc Cart API
//
// Documentation for Cart API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package doc

import "github.com/fedo3nik/cart-go-api/internal/domain/model"

// swagger:parameters addItemParams addItem
type addItemParams struct {
	// in: path
	// example: 1
	CartID int
	// in: body
	// example: "Hat"
	Product string `json:"product"`
	// in: body
	// example: 10
	Quantity int `json:"quantity"`
}

// swagger:parameters removeItemParams removeItem
type removeItemParams struct {
	// in: path
	// example: 3
	CartID int
	// in: path
	// example: 5
	ItemID int
}

// swagger:parameters getCartParams getCart
type getCartParams struct {
	// in: path
	// example: 1
	CartID int
}

// New cart created successfully
// swagger:response createCartResponse
type createCartResponse struct {
	// ID of the new cart
	CartID int `json:"id"`
	// Empty array of cartItems
	Items []model.CartItem `json:"items"`
}

// CartItem added into the cart successfully
// swagger:response addItemResponse
type addItemResponse struct {
	// ID of the new cartItem
	ID int `json:"id"`
	// CartID in which item was placed
	CartID int `json:"cart_id"`
	// Product title
	Product string `json:"product"`
	// Quantity of the products in the cartItem
	Quantity int `json:"quantity"`
}

// CartItem removed from the cart successfully
// swagger:response removeItemResponse
type removeItemResponse struct {
}

// The cart with the items in it
// swagger:response getCartResponse
type getCartResponse struct {
	// ID of the new cart
	CartID int `json:"id"`
	// Array of items placed in the cart
	Items []model.CartItem `json:"items"`
}

// Error caused
// swagger:response errorResponse
type errorResponse struct {
	// Error message
	Message string `json:"message"`
}
