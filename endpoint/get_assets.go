package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

// GetAssets endpoint
func GetAssets(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get date
		untilDate := r.URL.Query().Get("untildate")
		if untilDate == "" {
			untilDate = time.Now().Format(helper.DateLayoutYMD)
		} else {
			// check time validity
			if _, err := time.Parse(helper.DateLayoutYMD, untilDate); err != nil {
				WriteErrors(w, FieldErrors{{"date value", ErrFormatInvalid}})
				return
			}
		}

		// get assets data
		assets, err := app.Store.GetGoodsAssets(untilDate)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve assets", ErrFailed}})
			return
		}

		// reformat
		var details []model.AssetDetail
		var totalQty, totalPrice int
		for _, data := range assets {
			variantID, _ := strconv.Atoi(data[0])
			qty, _ := strconv.Atoi(data[4])
			averagePrice, _ := strconv.Atoi(data[5])
			total := qty * averagePrice
			details = append(details, model.AssetDetail{
				ProductVariantID: variantID,
				SKU:              data[1],
				Name:             fmt.Sprintf("%s (%s)", data[2], data[3]),
				AveragePrice:     averagePrice,
				Qty:              qty,
				Total:            total,
			})

			totalQty += qty
			totalPrice += total
		}

		result := model.Asset{
			Date:                 untilDate,
			TotalProductVariants: len(assets),
			TotalQty:             totalQty,
			TotalPrice:           totalPrice,
			Details:              details,
		}

		WriteData(w, http.StatusOK, result)
	}
}
