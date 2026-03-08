package domain

type Variant struct {
	ID        int64
	ProductID int64
	SKU       string
	Price     int64
	Stock     int
	Product   Product
}
