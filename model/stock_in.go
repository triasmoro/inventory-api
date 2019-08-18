package model

type StockIn struct {
	ID                    int    `json:"id,omitempty"`
	PurchaseOrderDetailID int    `json:"purchase_order_detail_id"`
	Time                  string `json:"time"`
	ReceiveQty            int    `json:"receive_qty"`
}
