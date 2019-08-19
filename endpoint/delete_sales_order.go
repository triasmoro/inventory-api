package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// DeleteSalesOrder endpoint
func DeleteSalesOrder(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		store := app.Store

		// check stock-in from this po
		exist, err := store.IsStockOutExistBySalesOrderID(id)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve stock-out", ErrFailed}})
			return
		} else if exist {
			WriteErrors(w, FieldErrors{{"stock-out exist", ErrFailed}})
			return
		}

		if err := store.DeleteSalesOrder(id); err != nil {
			WriteErrors(w, FieldErrors{{"delete sales order", ErrFailed}})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
