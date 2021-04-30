package service

import e "github.com/fedo3nik/cart-go-api/internal/errors"

// ValidateItemData validate product tittle and quantity of products in a new item.
// Returns ErrInvalidProduct error if product title is blank and
// returns ErrInvalidQuantity error if quantity of products less than 1.
// In other cases returns nil.
func (c CartService) ValidateItemData(product string, quantity int) error {
	if product == "" {
		return e.ErrInvalidProduct
	}

	if quantity <= 0 {
		return e.ErrInvalidQuantity
	}

	return nil
}
