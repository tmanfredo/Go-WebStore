package types

// TODO: If you choose to use a struct rather than individual parameters to your view, you might flesh this one out:
type PurchaseInfo struct {
	First string
	Last string
	Email string
	Product string
	Price float64
	Quantity int
	Donate string
	Tax float64
	Subtotal float64
	Total float64
}
