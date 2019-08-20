package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// DeleteStockOut endpoint
func DeleteStockOut(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		if err := app.Store.DeleteStockOut(id); err != nil {
			WriteErrors(w, FieldErrors{{"delete stock-out", ErrFailed}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
