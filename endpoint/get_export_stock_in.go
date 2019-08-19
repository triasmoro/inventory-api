package endpoint

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/helper"
)

const waitingNote = "Masih menunggu"

// GetExportStockIn endpoint
func GetExportStockIn(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// range date
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")
		if startDate == "" || endDate == "" {
			WriteErrors(w, FieldErrors{{"query parameter", ErrRequired}})
			return
		}

		// check time validity
		if _, err := time.Parse(helper.DateLayoutYMD, startDate); err != nil {
			WriteErrors(w, FieldErrors{{"start date", ErrFormatInvalid}})
			return
		}
		if _, err := time.Parse(helper.DateLayoutYMD, endDate); err != nil {
			WriteErrors(w, FieldErrors{{"end date", ErrFormatInvalid}})
			return
		}

		// get data
		stocks, err := app.Store.GetStockInReport(startDate, endDate)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve stock-in", ErrFailed}})
			return
		}

		buf := &bytes.Buffer{}
		writer := csv.NewWriter(buf)
		writer.Comma = ';'

		// preparing data
		records := rearrange(stocks)

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

		filename := fmt.Sprintf("stockin-%s.csv", time.Now().Format(helper.DateLayoutPlainYMDHIS))

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		w.Header().Set("Access-Control-Expose-Headers", "X-Suggested-Filename")
		w.Header().Set("X-Suggested-Filename", filename)
		w.Write(buf.Bytes())
	}
}

func rearrange(raw [][]string) [][]string {
	records := [][]string{
		// header field
		{
			"purchase time",
			"sku",
			"product name",
			"purchased qty",
			"received qty",
			"purchase price",
			"total",
			"po number",
			"notes",
		},
	}

	for _, data := range raw {

		// sum received qty
		receivedQty := 0
		for _, str := range strings.Split(data[10], ",") {
			qty, _ := strconv.Atoi(str)
			receivedQty += qty
		}

		// po number
		poNumber := "(Hilang)"
		if data[9] != "" {
			poNumber = data[9]
		}

		purchasedQty, _ := strconv.Atoi(data[6])

		records = append(records, []string{
			data[2],
			data[3],
			fmt.Sprintf("%s (%s)", data[4], data[5]),
			data[6],
			strconv.Itoa(receivedQty),
			data[7],
			data[8],
			poNumber,
			notes(purchasedQty, data[10], data[11]),
		})
	}

	return records
}

// generate notes
func notes(purchasedQty int, receivedQtys, receivedTimes string) string {
	qtys := strings.Split(receivedQtys, ",")
	times := strings.Split(receivedTimes, ",")

	var notes []string
	totalQty := 0
	for i, str := range qtys {
		qty, _ := strconv.Atoi(str)
		totalQty += qty

		notes = append(notes, fmt.Sprintf("%s terima %d", times[i], qty))
	}

	// add notes "Waiting" if there is pending goods
	if totalQty < purchasedQty {
		notes = append(notes, waitingNote)
	}

	return strings.Join(notes, "; ")
}
