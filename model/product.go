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
