package service

import (
	"context"

	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	"github.com/fedo3nik/cart-go-api/internal/infrastructure/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Cart interface {
	CreateCart(ctx context.Context) (*models.Cart, error)
	AddItem(ctx context.Context, product string, quantity, cartID int) (*models.CartItem, error)
	RemoveItem(ctx context.Context, cartID, itemID int) error
	GetCart(ctx context.Context, cartID int) (*models.Cart, error)
}

type CartService struct {
	Pool *pgxpool.Pool
}

func (c CartService) CreateCart(ctx context.Context) (*models.Cart, error) {
	var cart models.Cart

	id, err := postgres.InsertCart(ctx, c.Pool)
	if err != nil {
		return nil, e.ErrDB
	}

	cart.ID = id

	return &cart, nil
}

func (c CartService) AddItem(ctx context.Context, product string, quantity, cartID int) (*models.CartItem, error) {
	err := c.ValidateItemData(product, quantity)
	if err != nil {
		return nil, err
	}

	item := models.CartItem{Product: product, Quantity: quantity, CartID: cartID}

	id, err := postgres.InsertItem(ctx, c.Pool, &item)
	if err != nil {
		return nil, e.ErrInvalidCartID
	}

	item.ID = id

	return &item, nil
}

func (c CartService) RemoveItem(ctx context.Context, cartID, itemID int) error {
	flag, err := postgres.DeleteItem(ctx, c.Pool, cartID, itemID)
	if err != nil {
		return e.ErrDB
	}

	if !flag {
		return e.ErrRemove
	}

	return nil
}

func (c CartService) GetCart(ctx context.Context, cartID int) (*models.Cart, error) {
	cart, err := postgres.GetCart(ctx, c.Pool, cartID)
	if err != nil {
		return nil, e.ErrDB
	}

	if cart.ID == -1 {
		return nil, e.ErrInvalidCartID
	}

	return cart, nil
}

func NewCartService(pool *pgxpool.Pool) *CartService {
	return &CartService{Pool: pool}
}
