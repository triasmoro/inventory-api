package storage

import (
	"database/sql"
	"fmt"

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

// DeleteStockIn method
func (s *Store) DeleteStockIn(id int) error {
	query := fmt.Sprintf(`UPDATE stock_in SET fg_delete = 1 WHERE id = %d`, id)
	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

// DeleteStockOut method
func (s *Store) DeleteStockOut(id int) error {
	query := fmt.Sprintf(`UPDATE stock_out SET fg_delete = 1 WHERE id = %d`, id)
	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

// GetActualStocks method
func (s *Store) GetActualStocks() ([][]string, error) {
	query := `SELECT
		pv.id AS variant_id,
		pv.sku,
		p.name,
		GROUP_CONCAT(DISTINCT(pov.value)) AS options,
		COALESCE(SUM(st_in.receive_qty), 0) AS total_in,
		COALESCE(SUM(st_out.qty), 0) AS total_out
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
	var variantID, sku, name, options, totalIn, totalOut string
	for rows.Next() {
		err = rows.Scan(
			&variantID,
			&sku,
			&name,
			&options,
			&totalIn,
			&totalOut,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, []string{
			variantID,
			sku,
			name,
			options,
			totalIn,
			totalOut,
		})
	}

	return result, nil
}

// GetGoodsAssets method
func (s *Store) GetGoodsAssets(untilDate string) ([][]string, error) {
	query := `SELECT
		pv.id AS variant_id,
		pv.sku,
		p.name,
		GROUP_CONCAT(DISTINCT(pov.value)) AS options,
		COALESCE(SUM(pod.price), 0) / COALESCE(SUM(pod.qty), 1) AS average_price,
		COALESCE(SUM(st_in.receive_qty), 0) - COALESCE(SUM(st_out.qty), 0) AS stock
	FROM product_variants pv
	INNER JOIN products p ON p.id = pv.product_id
	INNER JOIN product_variant_options pvo ON pvo.product_variant_id = pv.id
	INNER JOIN product_option_values pov ON pov.id = pvo.product_option_value_id
	LEFT JOIN purchase_order_details pod ON pod.product_variant_id = pv.id
	LEFT JOIN stock_in st_in ON st_in.purchase_order_detail_id = pod.id
		AND st_in.fg_delete = 0
		AND DATE(st_in.time) <= "?"
	LEFT JOIN stock_out st_out ON st_out.product_variant_id = pv.id
		AND st_out.fg_delete = 0
		AND DATE(st_out.time) <= "?"
	WHERE pv.fg_delete = 0
	GROUP BY p.id, pv.id, pvo.product_variant_id`

	rows, err := s.DB.Query(query, untilDate, untilDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]string
	var variantID, sku, name, options, averagePrice, stock string
	for rows.Next() {
		err = rows.Scan(
			&variantID,
			&sku,
			&name,
			&options,
			&averagePrice,
			&stock,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, []string{
			variantID,
			sku,
			name,
			options,
			averagePrice,
			stock,
		})
	}

	return result, nil
}

// IsStockInExistByPurchaseOrderID method
func (s *Store) IsStockInExistByPurchaseOrderID(id int) (bool, error) {
	var exist int
	query := `SELECT 1 FROM purchase_orders po
		INNER JOIN purchase_order_details pod ON pod.purchase_order_id = po.id
		INNER JOIN stock_in st_in ON st_in.purchase_order_detail_id = pod.id
		WHERE po.fg_delete = 0
		AND st_in.fg_delete = 0
		AND po.id = ?`
	if err := s.DB.QueryRow(query, id).Scan(&exist); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
