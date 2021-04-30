package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	"github.com/fedo3nik/cart-go-api/internal/domain/models"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	dto "github.com/fedo3nik/cart-go-api/internal/interface/controller/dtohttp"

	"github.com/gorilla/mux"
)

// A HTTPCreateCartHandler represents handler for CreateCart endpoint.
type HTTPCreateCartHandler struct {
	cartService service.Cart
}

// A HTTPAddItemHandler represents handler for AddItem endpoint.
type HTTPAddItemHandler struct {
	cartService service.Cart
}

// A HTTPRemoveItemHandler represents handler for RemoveItem endpoint.
type HTTPRemoveItemHandler struct {
	cartService service.Cart
}

// A HTTPGetCartHandler represents handler for GetCart endpoint.
type HTTPGetCartHandler struct {
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

	if errors.Is(err, e.ErrRemove) {
		w.WriteHeader(http.StatusBadRequest)

		return &dto.ErrorResponse{Message: "Cart or item with these IDs does not exist"}
	}

	return nil
}

func NewHTTPCreateCartHandler(cartService service.Cart) *HTTPCreateCartHandler {
	return &HTTPCreateCartHandler{cartService: cartService}
}

// ServeHTTP is a method to handle CreateCart endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// For creating a Cart model used method CreateCart from the service layer.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPCreateCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.CartResponse

	cart, err := hh.cartService.CreateCart(r.Context())
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			return
		}

		return
	}

	resp.ID = cart.ID
	resp.Items = []models.CartItem{}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHTTPAddItemHandler(cartService service.Cart) *HTTPAddItemHandler {
	return &HTTPAddItemHandler{cartService: cartService}
}

// ServeHTTP is a method to handle AddItem endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// Data about the item received from the Request body using json.Decode(),
// cartID received from the URL via func Vars() from the mux package.
// For creating CartItem model used method AddItem from the service layer.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPAddItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cartID, err := strconv.Atoi(mux.Vars(r)["cartID"])
	if err != nil {
		return
	}

	var req dto.AddItemRequest

	var resp dto.AddItemResponse

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	item, err := hh.cartService.AddItem(r.Context(), req.Product, req.Quantity, cartID)
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHTTPRemoveItemHandler(cartService service.Cart) *HTTPRemoveItemHandler {
	return &HTTPRemoveItemHandler{cartService: cartService}
}

// ServeHTTP is a method to handle RemoveItem endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// ItemID and cartID received from the URL via func Vars() from the mux package.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPRemoveItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cartID, err := strconv.Atoi(mux.Vars(r)["cartID"])
	if err != nil {
		return
	}

	itemID, err := strconv.Atoi(mux.Vars(r)["itemID"])
	if err != nil {
		return
	}

	var resp dto.RemoveItemResponse

	err = hh.cartService.RemoveItem(r.Context(), cartID, itemID)
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func NewHTTPGetCartHandler(cartService service.Cart) *HTTPGetCartHandler {
	return &HTTPGetCartHandler{cartService: cartService}
}

// ServeHTTP is a method to handle GetCart endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// cartID received from the URL via func Vars() from the mux package.
// Method GetCart used for received all the items from this cart.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPGetCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cartID, err := strconv.Atoi(mux.Vars(r)["cartID"])
	if err != nil {
		return
	}

	var resp dto.CartResponse

	cart, err := hh.cartService.GetCart(r.Context(), cartID)
	if err != nil {
		resp := handleError(w, err)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			return
		}

		return
	}

	resp.ID = cartID
	resp.Items = cart.Items

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
