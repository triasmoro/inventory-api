package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/triasmoro/inventory-api/helper"

	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/model"
)

// GetSalesReport endpoint
func GetSalesReport(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get date
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")
		if startDate == "" || endDate == "" {
			WriteErrors(w, FieldErrors{{"start date & end date", ErrRequired}})
			return
		}

		// check time validity
		startDateParsed, err := time.Parse(helper.DateLayoutYMD, startDate)
		if err != nil {
			WriteErrors(w, FieldErrors{{"start date", ErrFormatInvalid}})
			return
		}
		endDateParsed, err := time.Parse(helper.DateLayoutYMD, endDate)
		if err != nil {
			WriteErrors(w, FieldErrors{{"end date", ErrFormatInvalid}})
			return
		}

		// get sales report data
		raw, err := app.Store.GetSalesReport(startDate, endDate)
		if err != nil {
			WriteErrors(w, FieldErrors{{"retrieve sales report", ErrFailed}})
			return
		}

		log.Println(raw)

		// reformat
		var report model.SalesReport
		var prevDetail model.SalesReportDetail
		var details []model.SalesReportDetail
		var variants []model.SalesReportProductVariant

		var totalIncome, totalGrossProfit, totalSales, totalProducts int

		for _, data := range raw {
			salesOrderID, _ := strconv.Atoi(data[0])

			soTime, err := time.Parse(helper.DateLayoutYMDTHISZ, data[2])
			if err != nil {
				log.Println(err)
			}

			detail := model.SalesReportDetail{
				SalesOrderID: salesOrderID,
				SONumber:     data[1],
				Time:         soTime.Format(helper.DateLayoutYMDHIS),
			}

			if prevDetail.SalesOrderID != 0 && prevDetail.SalesOrderID != salesOrderID {
				totalSales++
				prevDetail.Variants = variants
				variants = []model.SalesReportProductVariant{} // reset
			}

			productVariantID, _ := strconv.Atoi(data[4])
			qty, _ := strconv.Atoi(data[8])
			sellingPrice, _ := strconv.Atoi(data[9])
			totalSellingPrice, _ := strconv.Atoi(data[10])
			purchasePrice, _ := strconv.Atoi(data[11])
			grossProfit := totalSellingPrice - (qty * purchasePrice)

			totalProducts++

			variants = append(variants, model.SalesReportProductVariant{
				ProductVariantID:  productVariantID,
				SKU:               data[5],
				Name:              fmt.Sprintf("%s (%s)", data[6], data[7]),
				Qty:               qty,
				SellingPrice:      sellingPrice,
				TotalSellingPrice: totalSellingPrice,
				PurchasePrice:     purchasePrice,
				Profit:            grossProfit,
			})

			totalIncome += totalSellingPrice
			totalGrossProfit += grossProfit
			prevDetail = detail
		}

		// insert last iteration data
		totalSales++
		prevDetail.Variants = variants
		details = append(details, prevDetail)
		report.Details = details
		report.TotalIncome = totalIncome
		report.TotalGrossProfit = totalGrossProfit
		report.TotalSales = totalSales
		report.TotalProducts = totalProducts
		report.Date = time.Now().Format("02 January 2006")
		report.PeriodStart = startDateParsed.Format("02 January 2006")
		report.PeriodEnd = endDateParsed.Format("02 January 2006")

		WriteData(w, http.StatusOK, report)
	}
}
