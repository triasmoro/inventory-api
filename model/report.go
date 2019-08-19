package model

type ActualStock struct {
	ProductID        int    `json:"product_id"`
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

type SalesReport struct {
	Date             string              `json:"date"`
	PeriodStart      string              `json:"period_start"`
	PeriodEnd        string              `json:"period_end"`
	TotalIncome      int                 `json:"total_income"`
	TotalGrossProfit int                 `json:"total_gross_profit"`
	TotalSales       int                 `json:"total_sales"`
	TotalProducts    int                 `json:"total_products"`
	Details          []SalesReportDetail `json:"details"`
}

type SalesReportDetail struct {
	SalesOrderID int                         `json:"sales_order_id"`
	SONumber     string                      `json:"so_number"`
	Time         string                      `json:"time"`
	Variants     []SalesReportProductVariant `json:"variants"`
}

type SalesReportProductVariant struct {
	ProductVariantID  int    `json:"product_variant_id"`
	SKU               string `json:"sku"`
	Name              string `json:"name"`
	Qty               int    `json:"qty"`
	SellingPrice      int    `json:"selling_price"`
	TotalSellingPrice int    `json:"total_selling_price"`
	PurchasePrice     int    `json:"purchase_price"`
	Profit            int    `json:"profit"`
}
