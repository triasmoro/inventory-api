package storage

import (
	"strconv"

	"github.com/triasmoro/inventory-api/model"
)

// SaveStockIn method
func (s *Store) SaveStockIn(stock *model.StockIn) error {
	stmt, err := s.DB.Prepare(`INSERT INTO stock_in
		(purchase_order_detail_id, time, receive_qty) VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		stock.PurchaseOrderDetailID,
		stock.Time,
		stock.ReceiveQty,
	)
	if err != nil {
		return err
	}

	stockID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	stock.ID = int(stockID)

	return nil
}

// SaveStockOutWithSales method
func (s *Store) SaveStockOutWithSales(stock *model.StockOut) error {
	stmt, err := s.DB.Prepare(`INSERT INTO stock_out
		(sales_order_detail_id, product_variant_id, time, qty, notes) VALUES
		(?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		stock.SalesOrderDetailID,
		stock.ProductVariantID,
		stock.Time,
		stock.Qty,
		stock.Notes,
	)
	if err != nil {
		return err
	}

	stockID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	stock.ID = int(stockID)

	return nil
}

// SaveStockOutWithoutSales method
func (s *Store) SaveStockOutWithoutSales(stock *model.StockOut) error {
	stmt, err := s.DB.Prepare(`INSERT INTO stock_out
		(product_variant_id, time, qty, notes) VALUES
		(?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		stock.ProductVariantID,
		stock.Time,
		stock.Qty,
		stock.Notes,
	)
	if err != nil {
		return err
	}

	stockID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	stock.ID = int(stockID)

	return nil
}

// GetActualStocks method
func (s *Store) GetActualStocks() ([][]string, error) {
	query := `SELECT
		pv.id AS variant_id,
		pv.sku,
		p.name,
		GROUP_CONCAT(DISTINCT(pov.value)) AS options,
		SUM(st_in.receive_qty) AS total_in,
		SUM(st_out.qty) AS total_out
	FROM product_variants pv
	INNER JOIN products p ON p.id = pv.product_id
	INNER JOIN product_variant_options pvo ON pvo.product_variant_id = pv.id
	INNER JOIN product_option_values pov ON pov.id = pvo.product_option_value_id
	LEFT JOIN purchase_order_details pod ON pod.product_variant_id = pv.id
	LEFT JOIN stock_in st_in ON st_in.purchase_order_detail_id = pod.id AND st_in.fg_delete = 0
	LEFT JOIN stock_out st_out ON st_out.product_variant_id = pv.id AND st_out.fg_delete = 0
	WHERE pv.fg_delete = 0
	GROUP BY p.id, pv.id, pvo.product_variant_id`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]string
	var variantID, sku, name, options string
	var totalInText, totalOutText *string
	for rows.Next() {
		err = rows.Scan(
			&variantID,
			&sku,
			&name,
			&options,
			&totalInText,
			&totalOutText,
		)
		if err != nil {
			return nil, err
		}

		var totalIn int
		if totalInText != nil {
			totalIn, _ = strconv.Atoi(*totalInText)
		}

		var totalOut int
		if totalOutText != nil {
			totalOut, _ = strconv.Atoi(*totalOutText)
		}

		result = append(result, []string{
			variantID,
			sku,
			name,
			options,
			strconv.Itoa(totalIn),
			strconv.Itoa(totalOut),
		})
	}

	return result, nil
}
