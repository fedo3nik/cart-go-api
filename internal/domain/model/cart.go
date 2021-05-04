package model

// Cart represents shopping cart in the online store.
type Cart struct {
	ID    int // ID of the cart
	Items []CartItem
}
