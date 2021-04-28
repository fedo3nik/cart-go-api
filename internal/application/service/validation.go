package service

import e "github.com/fedo3nik/cart-go-api/internal/errors"

func (c CartService) ValidateItemData(product string, quantity int) error {
	if product == "" {
		return e.ErrInvalidProduct
	}

	if quantity <= 0 {
		return e.ErrInvalidQuantity
	}

	return nil
}
