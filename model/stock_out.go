package model

type StockOut struct {
	ID                 int    `json:"id,omitempty"`
	SalesOrderDetailID int    `json:"sales_order_detail_id"`
	ProductVariantID   int    `json:"product_variant_id"`
	Time               string `json:"time"`
	Qty                int    `json:"qty"`
	Notes              string `json:"notes"`
}
