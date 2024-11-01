package types

// TODO: If you choose to use a struct rather than individual parameters to your view, you might flesh this one out:
type PurchaseInfo struct {
	Welcome string
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

type Product struct{
	Id int
	Name string
	Image string
	Price float64
	Instock int
}

type Order struct{
    Product_Name string
    Customer_Name string
    Quantity int
    Price float64 
    Tax float64
    Donation float64
    Timestamp int
}

type Customer struct{
	Id int
	First string
	Last string
	Email string
}