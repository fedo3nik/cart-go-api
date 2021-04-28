package service

import (
	"context"
	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	"github.com/fedo3nik/cart-go-api/internal/infrastructure/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	e "github.com/fedo3nik/cart-go-api/internal/errors"
)

type Cart interface {
	CreateCart(ctx context.Context) (*models.Cart, error)
}

type CartService struct {
	Pool *pgxpool.Pool
}

func (c CartService) CreateCart(ctx context.Context) (*models.Cart, error) {
	var cart models.Cart

	id, err := postgres.Insert(ctx, c.Pool)
	if err != nil {
		return nil, errors.Wrap(e.ErrDB, "insert")
	}

	cart.ID = id

	return &cart, nil
}

func NewCartService(pool *pgxpool.Pool) *CartService {
	return &CartService{Pool: pool}
}
