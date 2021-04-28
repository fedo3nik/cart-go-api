package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	"github.com/gorilla/mux"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	dto "github.com/fedo3nik/cart-go-api/internal/interface/controller/dtohttp"
)

type HTTPCreateCartHandler struct {
	cartService service.Cart
}

type HTTPAddItemHandler struct {
	cartService service.Cart
}

func handleError(w http.ResponseWriter, err error) *dto.ErrorResponse {
	if errors.Is(err, e.ErrDB) {
		w.WriteHeader(http.StatusBadGateway)

		return &dto.ErrorResponse{Message: "Database error"}
	}

	if errors.Is(err, e.ErrInvalidCartID) {
		w.WriteHeader(http.StatusBadRequest)

		return &dto.ErrorResponse{Message: "Cart with the same ID does not exist"}
	}

	if errors.Is(err, e.ErrInvalidQuantity) {
		w.WriteHeader(http.StatusBadRequest)

		return &dto.ErrorResponse{Message: "Products quantity must be positive"}
	}

	if errors.Is(err, e.ErrInvalidProduct) {
		w.WriteHeader(http.StatusBadRequest)

		return &dto.ErrorResponse{Message: "Product title can't be blank"}
	}

	return nil
}

func NewHTTPCreateCartHandler(cartService service.Cart) *HTTPCreateCartHandler {
	return &HTTPCreateCartHandler{cartService: cartService}
}

func (hh HTTPCreateCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.CreateCartResponse

	cart, err := hh.cartService.CreateCart(r.Context())
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			log.Printf("Encode error: %v", err)

			return
		}

		return
	}

	resp.ID = cart.ID
	resp.Items = []models.CartItem{}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHTTPAddItemHandler(cartService service.Cart) *HTTPAddItemHandler {
	return &HTTPAddItemHandler{cartService: cartService}
}

func (hh HTTPAddItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cartID, err := strconv.Atoi(mux.Vars(r)["cartID"])
	if err != nil {
		log.Printf("Strconv err: %v", err)
		return
	}

	var req dto.AddItemRequest

	var resp dto.AddItemResponse

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad request: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	item, err := hh.cartService.AddItem(r.Context(), req.Product, req.Quantity, cartID)
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			log.Printf("Encode error: %v", err)

			return
		}

		return
	}

	resp.ID = item.ID
	resp.CartID = item.CartID
	resp.Product = item.Product
	resp.Quantity = item.Quantity

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
