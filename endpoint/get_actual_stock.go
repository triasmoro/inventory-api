package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/triasmoro/inventory-api/app"
)

// GetActualStock endpoint
func GetActualStock(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get stock data
		stocks, err := app.Store.GetActualStocks()
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve actual stocks", ErrFailed}})
			return
		}

		type ActualStock struct {
			ProductVariantID int    `json:"product_variant_id"`
			SKU              string `json:"sku"`
			Name             string `json:"product_name"`
			StockIn          int    `json:"total_stock_in"`
			StockOut         int    `json:"total_stock_out"`
		}

		// reformat
		var result []ActualStock
		for _, data := range stocks {
			variantID, _ := strconv.Atoi(data[0])
			stockIn, _ := strconv.Atoi(data[4])
			stockOut, _ := strconv.Atoi(data[5])
			result = append(result, ActualStock{
				ProductVariantID: variantID,
				SKU:              data[1],
				Name:             fmt.Sprintf("%s (%s)", data[2], data[3]),
				StockIn:          stockIn,
				StockOut:         stockOut,
			})
		}

		WriteData(w, http.StatusOK, result)
	}
}
