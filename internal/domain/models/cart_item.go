package models

type CartItem struct {
	ID       int
	CartID   int
	Quantity int
	Product  string
}
