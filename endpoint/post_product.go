package endpoint

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
	"github.com/triasmoro/inventory-api/model"
)

// PostProduct to add new product along with its variants
func PostProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var product model.Product
		if err := json.Unmarshal(body, &product); err != nil {
			WriteErrors(w, FieldErrors{{"unmarshal body", ErrFormatInvalid}})
			return
		}

		store := app.Store

		// check variant options
		for i, variant := range product.Variants {
			sku, err := helper.GenerateSKU()
			if err != nil {
				WriteErrors(w, FieldErrors{{"generate sku", ErrFailed}})
				return
			}

			// assign new sku to variant
			product.Variants[i].SKU = sku

			// get option id
			for j, option := range variant.Options {
				optionID, err := store.GetOptionValueID(option.Name, option.Value)
				if err != nil {
					WriteErrors(w, FieldErrors{{"retrieve option value", ErrNotFound}})
					return
				}

				// assign option id
				product.Variants[i].Options[j].ID = optionID
			}
		}

		// save
		if err := store.SaveProduct(&product); err != nil {
			WriteErrors(w, FieldErrors{{"save product", ErrFailed}})
			return
		}

		WriteData(w, http.StatusOK, product)
	}
}
