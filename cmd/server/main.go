package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	"github.com/fedo3nik/cart-go-api/internal/config"
	"github.com/fedo3nik/cart-go-api/internal/interface/controller"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Process config file error: %v", err)
	}

	pool, err := pgxpool.Connect(context.Background(), c.PostgresURL)
	if err != nil {
		log.Fatalf("Connect to database error: %v", err)
	}

	handler := mux.NewRouter()

	cartService := service.NewCartService(pool)

	createCartHandler := controller.NewHTTPCreateCartHandler(cartService)
	addItemHandler := controller.NewHTTPAddItemHandler(cartService)
	removeItemHandler := controller.NewHTTPRemoveItemHandler(cartService)
	getCartHandler := controller.NewHTTPGetCartHandler(cartService)

	handler.Handle("/carts", createCartHandler).Methods(http.MethodPost)
	handler.Handle("/carts/{cartID}/items", addItemHandler).Methods(http.MethodPost)
	handler.Handle("/carts/{cartID}/items/{itemID}", removeItemHandler).Methods(http.MethodDelete)
	handler.Handle("/carts/{cartID}", getCartHandler)

	err = http.ListenAndServe(c.Host+c.Port, handler)
	if err != nil {
		log.Fatalf("Listen & Serve error: %v", err)
	}
}
