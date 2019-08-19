package model

// Product object
type Product struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name"`
	Variants []Variant `json:"variants"`
}

// Variant of product
type Variant struct {
	ID        int      `json:"id,omitempty"`
	ProductID int      `json:"product_id,omitempty"`
	SKU       string   `json:"sku,omitempty"`
	Options   []Option `json:"options"`
}

// Option is differentiator between variants
type Option struct {
	ID    int    `json:"id,omitempty"` // option value id
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ExportedProduct object
type ExportedProduct struct {
	ID       int               `json:"id,omitempty"`
	Name     string            `json:"name"`
	Variants []ExportedVariant `json:"variants"`
}

// ExportedVariant of product
type ExportedVariant struct {
	ID        int      `json:"id,omitempty"`
	ProductID int      `json:"product_id,omitempty"`
	SKU       string   `json:"sku,omitempty"`
	Options   []Option `json:"options"`
	StockIn   int      `json:"stock_in"`  // accumulated
	StockOut  int      `json:"stock_out"` // accumulated
}
