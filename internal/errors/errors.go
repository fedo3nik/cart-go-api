package errors

import "errors"

var ErrDB = errors.New("database error")
var ErrInvalidProduct = errors.New("products name must not be blank")
var ErrInvalidQuantity = errors.New("products quantity should be positive")
var ErrInvalidCartID = errors.New("cart with the same ID does not exist")
var ErrRemove = errors.New("cart or item with these IDs does not exist")
