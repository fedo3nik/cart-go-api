package service

import (
	"context"

	"github.com/fedo3nik/cart-go-api/internal/domain/model"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	"github.com/fedo3nik/cart-go-api/internal/infrastructure/database/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Cart is the interface that describes methods for the service layer.
type Cart interface {
	CreateCart(ctx context.Context) (*model.Cart, error)
	AddItem(ctx context.Context, product string, quantity, cartID int) (*model.CartItem, error)
	RemoveItem(ctx context.Context, cartID, itemID int) error
	GetCart(ctx context.Context, cartID int) (*model.Cart, error)
}

// CartService represents service layer.
type CartService struct {
	Pool *pgxpool.Pool // connection pool
}

// CreateCart creates a new cart.
// Returns a pointer to the cart model.
// Also it returns a database error if the used func InsertCart returns an error.
func (c CartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	var cart model.Cart

	id, err := postgres.InsertCart(ctx, c.Pool)
	if err != nil {
		return nil, errors.Wrap(e.ErrDB, err.Error())
	}

	cart.ID = id

	return &cart, nil
}

// AddItem adds a new item to the cart.
// Returns a pointer to the item model.
// Also it returns an error if the item data is invalid or
// the cart with the same id doesn't exist.
func (c CartService) AddItem(ctx context.Context, product string, quantity, cartID int) (*model.CartItem, error) {
	err := c.ValidateItemData(product, quantity)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	item := model.CartItem{Product: product, Quantity: quantity, CartID: cartID}

	id, err := postgres.InsertItem(ctx, c.Pool, &item)
	if err != nil {
		return nil, errors.Wrap(e.ErrInvalidCartID, err.Error())
	}

	item.ID = id

	return &item, nil
}

// RemoveItem removes item from the cart.
// Returns an error if cart or item with the received IDs doesn't exist.
func (c CartService) RemoveItem(ctx context.Context, cartID, itemID int) error {
	flag, err := postgres.DeleteItem(ctx, c.Pool, cartID, itemID)
	if err != nil {
		return errors.Wrap(e.ErrDB, err.Error())
	}

	if flag {
		return e.ErrRemove
	}

	return nil
}

// GetCart gets the data about the cart with the ID == cartID.
// Returns a pointer to the cart model.
// Also it returns an error if the cart with the same ID doesn't exist.
func (c CartService) GetCart(ctx context.Context, cartID int) (*model.Cart, error) {
	cart, err := postgres.GetCart(ctx, c.Pool, cartID)
	if err != nil {
		return nil, errors.Wrap(e.ErrDB, err.Error())
	}

	if cart.ID == -1 {
		return nil, e.ErrInvalidCartID
	}

	return cart, nil
}

// NewCartService is a constructor for CartService struct.
func NewCartService(pool *pgxpool.Pool) *CartService {
	return &CartService{Pool: pool}
}
