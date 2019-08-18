package model

type ActualStock struct {
	ProductVariantID int    `json:"product_variant_id"`
	SKU              string `json:"sku"`
	Name             string `json:"product_name"`
	StockIn          int    `json:"total_stock_in"`
	StockOut         int    `json:"total_stock_out"`
}

type Asset struct {
	Date                 string        `json:"date"`
	TotalProductVariants int           `json:"total_product_variants"`
	TotalQty             int           `json:"total_qty"`
	TotalPrice           int           `json:"total_price"`
	Details              []AssetDetail `json:"details"`
}

type AssetDetail struct {
	ProductVariantID int    `json:"product_variant_id"`
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Qty              int    `json:"qty"`
	AveragePrice     int    `json:"average_price"`
	Total            int    `json:"total"`
}
