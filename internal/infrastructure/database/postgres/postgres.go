package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func Insert(ctx context.Context, p *pgxpool.Pool) (int, error) {
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
