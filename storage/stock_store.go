package storage

import "github.com/triasmoro/inventory-api/model"

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
