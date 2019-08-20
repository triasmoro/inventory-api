package storage

import (
	"database/sql"
	"fmt"

	"github.com/triasmoro/inventory-api/model"
)

// SaveSalesOrder method
func (s *Store) SaveSalesOrder(order *model.SalesOrder) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// save order
	stmt, err := tx.Prepare("INSERT INTO sales_orders(so_number, time) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := stmt.Exec(order.SONumber, order.Time)
	if err != nil {
		tx.Rollback()
		return err
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// assign id generated
	order.ID = int(orderID)

	for i, detail := range order.Details {
		// save order detail
		stmt, err := tx.Prepare(`INSERT INTO sales_order_details
			(sales_order_id, product_variant_id, qty, price) VALUES
			(?, ?, ?, ?)`)
		if err != nil {
			tx.Rollback()
			return err
		}

		res, err := stmt.Exec(orderID, detail.ProductVariantID, detail.Qty, detail.Price)
		if err != nil {
			tx.Rollback()
			return err
		}

		orderDetailID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		// assign id generated
		order.Details[i].ID = int(orderDetailID)
	}

	tx.Commit()

	return nil
}

// DeleteSalesOrder method
func (s *Store) DeleteSalesOrder(id int) error {
	query := fmt.Sprintf(`UPDATE sales_orders SET fg_delete = 1 WHERE id = %d`, id)
	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

// GetSalesOrderDetailAndSONumberByID method
func (s *Store) GetSalesOrderDetailAndSONumberByID(id int) (model.SalesOrderDetail, string, error) {
	var orderDetail model.SalesOrderDetail
	var orderNumber string
	query := `SELECT
			sod.product_variant_id,
			sod.qty,
			sod.price,
			so.so_number
		FROM sales_orders so
		INNER JOIN sales_order_details sod ON sod.sales_order_id = so.id
		WHERE sod.id = ? AND so.fg_delete = 0`
	if err := s.DB.QueryRow(query, id).Scan(
		&orderDetail.ProductVariantID,
		&orderDetail.Qty,
		&orderDetail.Price,
		&orderNumber); err != nil {
		if err != sql.ErrNoRows {
			return orderDetail, "", err
		}

		return orderDetail, "", nil
	}

	return orderDetail, orderNumber, nil
}

// GetSalesReport method
func (s *Store) GetSalesReport(start, end string) ([][]string, error) {
	query := fmt.Sprintf(`SELECT
		so.id AS sales_order_id,
		so.so_number,
		so.time,
		sod.id AS sales_order_detail_id,
		pv.id AS product_variant_id,
		pv.sku,
		p.name,
		GROUP_CONCAT(DISTINCT(pov.value)) AS options,
		sod.qty,
		sod.price AS selling_price,
		sod.qty * sod.price AS total,
		COALESCE(SUM(pod.price), 0) / COALESCE(SUM(pod.qty), 1) AS average_price
	FROM sales_orders so
	INNER JOIN sales_order_details sod ON sod.sales_order_id = so.id
	INNER JOIN product_variants pv ON pv.id = sod.product_variant_id
	INNER JOIN products p ON p.id = pv.product_id
	INNER JOIN product_variant_options pvo ON pvo.product_variant_id = pv.id
	INNER JOIN product_option_values pov ON pov.id = pvo.product_option_value_id
	INNER JOIN purchase_order_details pod ON pod.product_variant_id = pv.id
	INNER JOIN purchase_orders po ON po.id = pod.purchase_order_id
	WHERE so.fg_delete = 0
	AND po.fg_delete = 0
	AND pv.fg_delete = 0
	AND DATE(so.time) BETWEEN "%s" AND "%s"
	GROUP BY pv.id
	ORDER BY so.id, sod.id, pv.id ASC`, start, end)

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]string
	var (
		salesOrderID,
		salesOrderNumber,
		soDate,
		salesOrderDetailID,
		productVariantID,
		sku,
		name,
		options,
		qty,
		sellingPrice,
		totalSellingPrice,
		averagePurchasePrice string
	)
	for rows.Next() {
		err = rows.Scan(
			&salesOrderID,
			&salesOrderNumber,
			&soDate,
			&salesOrderDetailID,
			&productVariantID,
			&sku,
			&name,
			&options,
			&qty,
			&sellingPrice,
			&totalSellingPrice,
			&averagePurchasePrice,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, []string{
			salesOrderID,
			salesOrderNumber,
			soDate,
			salesOrderDetailID,
			productVariantID,
			sku,
			name,
			options,
			qty,
			sellingPrice,
			totalSellingPrice,
			averagePurchasePrice,
		})
	}

	return result, nil
}
