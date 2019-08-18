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

// SaveStockOutBasedSalesOrder method
func (s *Store) SaveStockOutBasedSalesOrder(stock *model.StockOut) error {
	stmt, err := s.DB.Prepare(`INSERT INTO stock_out
		(sales_order_detail_id, time, qty, notes) VALUES
		(?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		stock.SalesOrderDetailID,
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

// SaveStockOutBasedProduct method
func (s *Store) SaveStockOutBasedProduct(stock *model.StockOut) error {
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
