package endpoint

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

// GetExportProduct endpoint
// Read this about csv separator https://stackoverflow.com/questions/10140999/csv-with-comma-or-semicolon
func GetExportProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get products
		products, err := app.Store.ExportProducts()
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve products", ErrFailed}})
			return
		}

		buf := &bytes.Buffer{}
		writer := csv.NewWriter(buf)
		writer.Comma = ';'

		// convert products object into plain slice string
		records := flat(products)

		// write
		if err := writer.WriteAll(records); err != nil {
			WriteErrors(w, FieldErrors{{"writing csv record", ErrFailed}})
			return
		}

		// Write any buffered data to the underlying writer (standard output).
		writer.Flush()

		if err := writer.Error(); err != nil {
			WriteErrors(w, FieldErrors{{"writing csv", ErrFailed}})
			return
		}

		filename := fmt.Sprintf("product-%s.csv", time.Now().Format(helper.DateLayoutPlainYMDHIS))

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		w.Header().Set("Access-Control-Expose-Headers", "X-Suggested-Filename")
		w.Header().Set("X-Suggested-Filename", filename)
		w.Write(buf.Bytes())
	}
}

func flat(products []model.Product) [][]string {
	records := [][]string{
		{"product name", "sku", "option name", "option value"}, // header field
	}
	var prevProductID, prevVariantID int

	for _, product := range products {
		for _, variant := range product.Variants {
			for _, option := range variant.Options {

				var productName, sku string
				if product.ID != prevProductID {
					productName = product.Name
				}
				if variant.ID != prevVariantID {
					sku = variant.SKU
				}

				record := []string{
					productName,
					sku,
					option.Name,
					option.Value,
				}
				records = append(records, record)

				prevProductID = product.ID
				prevVariantID = variant.ID
			}
		}
	}

	return records
}
