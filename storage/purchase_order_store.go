package storage

import "github.com/triasmoro/inventory-api/model"

// SavePurchaseOrder method
func (s *Store) SavePurchaseOrder(order *model.PurchaseOrder) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// save order
	stmt, err := tx.Prepare("INSERT INTO purchase_orders(po_number, time) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := stmt.Exec(order.PONumber, order.Time)
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
		stmt, err := tx.Prepare(`INSERT INTO purchase_order_details
			(purchase_order_id, product_variant_id, qty, price) VALUES
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
