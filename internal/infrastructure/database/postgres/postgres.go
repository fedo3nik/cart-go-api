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

func DeleteItem(ctx context.Context, p *pgxpool.Pool, cartID, itemID int) (bool, error) {
	conn, err := p.Acquire(ctx)
	if err != nil {
		return true, err
	}

	defer conn.Release()

	ct, err := conn.Exec(ctx, "DELETE FROM Items WHERE ID=$1 AND cartID=$2", itemID, cartID)
	if err != nil {
		return true, err
	}

	if ct.RowsAffected() != 1 {
		return false, nil
	}

	return true, nil
}

func GetCart(ctx context.Context, p *pgxpool.Pool, cartID int) (*models.Cart, error) {
	var items []models.CartItem

	var cart models.Cart

	conn, err := p.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT id, cartId, product_name, quantity FROM items WHERE cartID=$1", cartID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.CartItem

		err = rows.Scan(&item.ID, &item.CartID, &item.Product, &item.Quantity)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	cart.ID = cartID
	cart.Items = items

	return &cart, nil
}
