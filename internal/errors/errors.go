package errors

import "errors"

// ErrDB is a custom error that returns if the error cuddle with database.
var ErrDB = errors.New("database error")

// ErrInvalidProduct is a custom error that returns if product title is blank.
var ErrInvalidProduct = errors.New("products name must not be blank")

// ErrInvalidQuantity is a custom error that returns if number of products less than 1.
var ErrInvalidQuantity = errors.New("products quantity should be positive")

// ErrInvalidCartID is a custom error that returns if cart with the same ID doesn't exist.
var ErrInvalidCartID = errors.New("cart with the same ID does not exist")

// ErrRemove is a custom error that returns if user try to remove non-existent item or from non-existing cart.
var ErrRemove = errors.New("cart or item with these IDs does not exist")
