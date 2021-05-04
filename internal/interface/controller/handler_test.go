package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	"github.com/fedo3nik/cart-go-api/internal/config"
	"github.com/fedo3nik/cart-go-api/internal/domain/model"
	"github.com/fedo3nik/cart-go-api/internal/infrastructure/database/postgres"
	controller "github.com/fedo3nik/cart-go-api/internal/interface/controller/dtohttp"

	"github.com/gavv/httpexpect/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestHTTPAddItemHandler_ServeHTTP(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	var maxItemID int

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)

	err = conn.QueryRow(context.Background(), "SELECT MAX(id) FROM items").Scan(&maxItemID)
	require.NoError(t, err)

	cartService := service.NewCartService(pool)

	handler := NewHTTPAddItemHandler(cartService)

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	testItem := controller.AddItemRequest{Product: "test_product", Quantity: 1}

	e.POST("/carts/3/items").WithJSON(testItem).Expect().Status(http.StatusOK)
}

func TestHTTPCreateCartHandler_ServeHTTP(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	var maxCartID int

	conn, err := pool.Acquire(context.Background())
	require.NoError(t, err)

	err = conn.QueryRow(context.Background(), "SELECT MAX(id) FROM carts").Scan(&maxCartID)
	require.NoError(t, err)

	cartService := service.NewCartService(pool)

	handler := NewHTTPCreateCartHandler(cartService)

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	expectedResult := controller.CartResponse{ID: maxCartID + 1, Items: []model.CartItem{}}

	e.POST("/carts").Expect().Status(http.StatusOK).JSON().Object().Equal(expectedResult)
}

func TestHTTPGetCartHandler_ServeHTTP(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	cartService := service.NewCartService(pool)

	handler := NewHTTPGetCartHandler(cartService)

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/carts/1").Expect().Status(http.StatusOK).Body().Empty()
}

func TestHTTPRemoveItemHandler_ServeHTTP(t *testing.T) {
	c, err := config.NewConfig()
	require.NoError(t, err)

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	require.NoError(t, err)

	id, err := postgres.InsertItem(context.Background(), pool, &model.CartItem{CartID: 3})
	require.NoError(t, err)

	itemID := strconv.Itoa(id)

	cartService := service.NewCartService(pool)

	handler := NewHTTPRemoveItemHandler(cartService)

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.DELETE("/carts/3/items/" + itemID).Expect().Status(http.StatusOK)
}
