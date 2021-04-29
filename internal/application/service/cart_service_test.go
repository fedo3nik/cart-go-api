package service

import (
	"context"
	"log"
	"testing"

	"github.com/fedo3nik/cart-go-api/internal/config"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCartService_AddItem(t *testing.T) {
	tests := []struct {
		name           string
		product        string
		cartID         int
		quantity       int
		expectedResult string
		expectedError  error
	}{
		{
			name:           "Add valid item",
			cartID:         2,
			product:        "Test_product",
			quantity:       5,
			expectedResult: "Test_product",
			expectedError:  nil,
		},
		{
			name:           "Invalid product title",
			product:        "",
			cartID:         1,
			quantity:       1,
			expectedResult: "",
			expectedError:  e.ErrInvalidProduct,
		},
		{
			name:           "Invalid quantity",
			product:        "Shoes",
			cartID:         1,
			quantity:       -1,
			expectedResult: "",
			expectedError:  e.ErrInvalidQuantity,
		},
	}

	c, err := config.NewConfig()
	assert.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	if err != nil {
		log.Fatalf("Connect to database error: %v", err)
	}

	cs := NewCartService(pool)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			item, err := cs.AddItem(context.Background(), tt.product, tt.quantity, tt.cartID)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.Equal(t, tt.expectedResult, item.Product)
		})
	}
}

func TestCartService_RemoveItem(t *testing.T) {
	tests := []struct {
		name          string
		cartID        int
		itemID        int
		expectedError error
	}{
		{
			name:          "Invalid cartID",
			cartID:        -1,
			itemID:        10,
			expectedError: e.ErrRemove,
		},
		{
			name:          "Invalid itemID",
			cartID:        1,
			itemID:        -1,
			expectedError: e.ErrRemove,
		},
	}

	c, err := config.NewConfig()
	assert.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	if err != nil {
		log.Fatalf("Connect to database error: %v", err)
	}

	cs := NewCartService(pool)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := cs.RemoveItem(context.Background(), tt.cartID, tt.itemID)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
