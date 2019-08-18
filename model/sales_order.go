package model

// SalesOrder object
type SalesOrder struct {
	ID       int                `json:"id,omitempty"`
	Time     string             `json:"time"`
	SONumber string             `json:"so_number"`
	Details  []SalesOrderDetail `json:"details"`
}

// SalesOrderDetail is detail that contains goods, qty, and price of SO
type SalesOrderDetail struct {
	ID               int `json:"id,omitempty"`
	ProductVariantID int `json:"product_variant_id"`
	Qty              int `json:"qty"`
	Price            int `json:"price"`
}
