package model

// CartItem represents items added to cart.
type CartItem struct {
	ID       int    // ID of the item
	CartID   int    // Cart ID to which this item was added
	Quantity int    // Quantity of the items in the cart
	Product  string // The name of the product from which this item was composed
}
