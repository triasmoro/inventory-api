package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/model"
)

// PatchProduct endpoint to only change product name
func PatchProduct(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErrors(w, FieldErrors{{"read body", ErrFormatInvalid}})
			return
		}

		var params model.Product
		if err := json.Unmarshal(body, &params); err != nil {
			WriteErrors(w, FieldErrors{{"unmarshal body", ErrFormatInvalid}})
			return
		}

		store := app.Store

		// check existing of product
		product, err := store.GetProductByID(id)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve product", ErrNotFound}})
			return
		}

		// update
		product.Name = params.Name
		if err := store.UpdateProduct(product); err != nil {
			WriteErrors(w, FieldErrors{{"update product", ErrNotFound}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
