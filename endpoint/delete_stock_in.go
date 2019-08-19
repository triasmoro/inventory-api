package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// DeleteStockIn endpoint
func DeleteStockIn(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		if err := app.Store.DeleteStockIn(id); err != nil {
			WriteErrors(w, FieldErrors{{"delete stock-in", ErrFailed}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
