package controller

import (
	"encoding/json"
	"errors"
	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	"log"
	"net/http"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	dto "github.com/fedo3nik/cart-go-api/internal/interface/controller/dtohttp"
)

type HTTPCreateCartHandler struct {
	cartService service.Cart
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, e.ErrDB) {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func NewHTTPCreateCartHandler(cartService service.Cart) *HTTPCreateCartHandler {
	return &HTTPCreateCartHandler{cartService: cartService}
}

func (hh HTTPCreateCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.CreateCartResponse

	cart, err := hh.cartService.CreateCart(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	resp.ID = cart.ID
	resp.Items = []models.CartItem{}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}
