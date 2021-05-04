package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/fedo3nik/cart-go-api/internal/application/service"
	"github.com/fedo3nik/cart-go-api/internal/domain/model"
	e "github.com/fedo3nik/cart-go-api/internal/errors"
	dto "github.com/fedo3nik/cart-go-api/internal/interface/controller/dtohttp"

	"github.com/gorilla/mux"
)

// HTTPCreateCartHandler represents handler for CreateCart endpoint.
type HTTPCreateCartHandler struct {
	cartService service.Cart
}

// HTTPAddItemHandler represents handler for AddItem endpoint.
type HTTPAddItemHandler struct {
	cartService service.Cart
}

// HTTPRemoveItemHandler represents handler for RemoveItem endpoint.
type HTTPRemoveItemHandler struct {
	cartService service.Cart
}

// HTTPGetCartHandler represents handler for GetCart endpoint.
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

// NewHTTPCreateCartHandler is a constructor for HTTPCreateCartHandler struct.
func NewHTTPCreateCartHandler(cartService service.Cart) *HTTPCreateCartHandler {
	return &HTTPCreateCartHandler{cartService: cartService}
}

// swagger:route POST /carts carts createCart
// Returns a new cart
// responses:
//	200: createCartResponse
//  502: errorResponse

// ServeHTTP is a method to handle CreateCart endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// For creating a Cart model used method CreateCart from the service layer.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPCreateCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.CartResponse

	w.Header().Set("Content-Type", "application/json")

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
	resp.Items = []model.CartItem{}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

// NewHTTPAddItemHandler is a constructor for HTTPAddItemHandler struct.
func NewHTTPAddItemHandler(cartService service.Cart) *HTTPAddItemHandler {
	return &HTTPAddItemHandler{cartService: cartService}
}

// swagger:route POST /carts/{cartID}/items items addItem
// Returns a new cartItem
// responses:
//	200: addItemResponse
//	400: errorResponse
//	502: errorResponse

// ServeHTTP is a method to handle AddItem endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// Data about the item received from the Request body using json.Decode(),
// cartID received from the URL via func Vars() from the mux package.
// For creating CartItem model used method AddItem from the service layer.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPAddItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	strCartID := mux.Vars(r)["cartID"]

	cartID, err := strconv.Atoi(strCartID)
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

// NewHTTPRemoveItemHandler is a constructor for HTTPRemoveItemHandler struct.
func NewHTTPRemoveItemHandler(cartService service.Cart) *HTTPRemoveItemHandler {
	return &HTTPRemoveItemHandler{cartService: cartService}
}

// swagger:route DELETE /carts/{cartID}/items/{itemID} items removeItem
// Returns empty json Object
// responses:
//	200: removeItemResponse
//	400: errorResponse
//	502: errorResponse

// ServeHTTP is a method to handle RemoveItem endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// ItemID and cartID received from the URL via func Vars() from the mux package.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPRemoveItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	strCartID := mux.Vars(r)["cartID"]
	strItemID := mux.Vars(r)["itemID"]

	cartID, err := strconv.Atoi(strCartID)
	if err != nil {
		return
	}

	itemID, err := strconv.Atoi(strItemID)
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

// NewHTTPGetCartHandler is a constructor for HTTPGetCartHandler struct.
func NewHTTPGetCartHandler(cartService service.Cart) *HTTPGetCartHandler {
	return &HTTPGetCartHandler{cartService: cartService}
}

// swagger:route GET /carts/{cartID} carts getCart
// Returns cart with the items in it
// responses:
//	200: getCartResponse
//	400: errorResponse
//	502: errorResponse

// ServeHTTP is a method to handle GetCart endpoint.
// It uses ResponseWriter and pointer to the Request from the standard package http.
// cartID received from the URL via func Vars() from the mux package.
// Method GetCart used for received all the items from this cart.
// Response write to the ResponseWriter using json.Encode().
func (hh HTTPGetCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	strCartID := mux.Vars(r)["cartID"]

	cartID, err := strconv.Atoi(strCartID)
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
