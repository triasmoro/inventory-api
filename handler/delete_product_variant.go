package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// DeleteProductVariant endpoint
func DeleteProductVariant(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		if err := app.Store.DeleteProductVariant(id); err != nil {
			WriteErrors(w, FieldErrors{{"delete product variant", ErrFailed}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
