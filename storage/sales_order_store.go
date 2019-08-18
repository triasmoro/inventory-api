package storage

import (
	"database/sql"

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

// GetSONumberByOrderDetailID method
func (s *Store) GetSONumberByOrderDetailID(id int) (string, error) {
	var orderNumber string
	query := `SELECT so.so_number
		FROM sales_orders so
		INNER JOIN sales_order_details sod ON sod.sales_order_id = so.id
		WHERE sod.id = ? AND so.fg_delete = 0`
	if err := s.DB.QueryRow(query, id).Scan(&orderNumber); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		return "", nil
	}

	return orderNumber, nil
}
