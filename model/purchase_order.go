package model

// PurchaseOrder object
type PurchaseOrder struct {
	ID       int                   `json:"id,omitempty"`
	Time     string                `json:"time"`
	PONumber string                `json:"po_number"`
	Details  []PurchaseOrderDetail `json:"details"`
}

// PurchaseOrderDetail is detail that contains goods, qty, and price of PO
type PurchaseOrderDetail struct {
	ID               int `json:"id,omitempty"`
	ProductVariantID int `json:"product_variant_id"`
	Qty              int `json:"qty"`
	Price            int `json:"price"`
}
