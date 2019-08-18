package model

// Product object
type Product struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name"`
	Variants []Variant `json:"variants"`
}
