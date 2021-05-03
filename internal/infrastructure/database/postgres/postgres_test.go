package postgres

import (
	"context"
	"testing"

	"github.com/fedo3nik/cart-go-api/internal/config"
	"github.com/fedo3nik/cart-go-api/internal/domain/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteItem(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	item := model.CartItem{CartID: 1, Quantity: 1, Product: "test_product"}

	id, err := InsertItem(context.Background(), pool, &item)
	require.NoError(t, err)

	tt := []struct {
		name           string
		cartID         int
		itemID         int
		expectedResult bool
	}{
		{
			name:           "Delete item",
			cartID:         1,
			itemID:         id,
			expectedResult: false,
		},
		{
			name:           "Wrong itemID",
			cartID:         1,
			itemID:         -1,
			expectedResult: true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			flag, err := DeleteItem(context.Background(), pool, tc.cartID, tc.itemID)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, flag)
		})
	}
}

func TestGetCart(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	cartID, err := InsertCart(context.Background(), pool)
	require.NoError(t, err)

	tt := []struct {
		name           string
		cartID         int
		expectedResult *model.Cart
	}{
		{
			name:           "Get empty cart",
			cartID:         cartID,
			expectedResult: &model.Cart{ID: cartID, Items: []model.CartItem{}},
		},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			cart, err := GetCart(context.Background(), pool, tc.cartID)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, cart)
		})
	}
}

func TestInsertCart(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	var maxCartID int

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)

	err = conn.QueryRow(context.Background(), "SELECT MAX(id) FROM carts").Scan(&maxCartID)
	require.NoError(t, err)

	tt := []struct {
		name           string
		expectedResult int
	}{
		{
			name:           "Insert new cart",
			expectedResult: maxCartID + 1,
		},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			id, err := InsertCart(context.Background(), pool)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, id)
		})
	}
}

func TestInsertItem(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	var maxItemID, maxCartID int

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)

	err = conn.QueryRow(context.Background(), "SELECT MAX(id) FROM items").Scan(&maxItemID)
	require.NoError(t, err)

	err = conn.QueryRow(context.Background(), "SELECT MAX(id) FROM carts").Scan(&maxCartID)
	require.NoError(t, err)

	tt := []struct {
		name           string
		item           model.CartItem
		expectedResult int
	}{
		{
			name:           "Insert item",
			item:           model.CartItem{CartID: maxCartID},
			expectedResult: maxItemID + 2,
		},
	}
	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			id, err := InsertItem(context.Background(), pool, &tc.item)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, id)
		})
	}
}
