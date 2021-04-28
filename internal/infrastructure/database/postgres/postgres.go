package postgres

import (
	"context"
	"log"

	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertCart(ctx context.Context, p *pgxpool.Pool) (int, error) {
	var id int

	conn, err := p.Acquire(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Release()

	row := conn.QueryRow(ctx, "INSERT INTO carts DEFAULT VALUES RETURNING id")

	err = row.Scan(&id)
	if err != nil {
		log.Printf("Scan error: %v", err)
		return 0, err
	}

	return id, nil
}

func InsertItem(ctx context.Context, p *pgxpool.Pool, item *models.CartItem) (int, error) {
	var id int

	conn, err := p.Acquire(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Release()

	row := conn.QueryRow(ctx, "INSERT INTO items (cartID, product_name, quantity) VALUES ($1, $2, $3) RETURNING id",
		item.CartID, item.Product, item.Quantity)

	err = row.Scan(&id)
	if err != nil {
		log.Printf("Scan error: %v", err)
		return 0, err
	}

	return id, nil
}
