package postgres

import (
	"context"

	"github.com/fedo3nik/cart-go-api/internal/domain/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InsertCart inserts a new Cart in the DB.
// Returns the ID of a new cart.
// Also it returns an error if the connection from the connection pool doesn't acquired or
// if a new cart doesn't inserted in the table.
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
		return 0, err
	}

	return id, nil
}

// InsertItem inserts a new CartItem in the DB.
// Returns the ID of a new item.
// Also it returns an error if the connection from the connection pool doesn't acquire or
// if a new item doesn't inserted in the table.
func InsertItem(ctx context.Context, p *pgxpool.Pool, item *model.CartItem) (int, error) {
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
		return 0, err
	}

	return id, nil
}

// DeleteItem deletes a CartItem from the cart in the DB.
// Returns the bool value that flagged item was deleted or no.
// Also it returns an error if the connection from the connection pool doesn't acquire or
// if the item doesn't deleted from the table.
func DeleteItem(ctx context.Context, p *pgxpool.Pool, cartID, itemID int) (bool, error) {
	conn, err := p.Acquire(ctx)
	if err != nil {
		return false, err
	}

	defer conn.Release()

	ct, err := conn.Exec(ctx, "DELETE FROM Items WHERE ID=$1 AND cartID=$2", itemID, cartID)
	if err != nil {
		return false, err
	}

	if ct.RowsAffected() != 1 {
		return true, nil
	}

	return false, nil
}

// GetCart selects all the items in the Cart from the DB.
// Returns pointer to the Cart model with the data.
// Also it returns an error if the connection from the connection pool doesn't acquire or
// if the error occurred while reading rows or if the cart doesn't selected from the table.
func GetCart(ctx context.Context, p *pgxpool.Pool, cartID int) (*model.Cart, error) {
	var items []model.CartItem

	var cart model.Cart

	var rowsCount int

	conn, err := p.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM carts WHERE ID=$1", cartID).Scan(&rowsCount)
	if err != nil {
		return nil, err
	}

	if rowsCount <= 0 {
		return &model.Cart{ID: -1}, nil
	}

	rows, err := conn.Query(ctx, "SELECT id, cartId, product_name, quantity FROM items WHERE cartID=$1", cartID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item model.CartItem

		err = rows.Scan(&item.ID, &item.CartID, &item.Product, &item.Quantity)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	cart.ID = cartID
	if items != nil {
		cart.Items = items
	} else {
		cart.Items = []model.CartItem{}
	}

	return &cart, nil
}
