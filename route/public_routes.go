package route

import (
	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/handler"
)

// publicRoutes to register all public endpoints
func publicRoutes(app *app.App) []*HandledRoute {
	var routes []*HandledRoute
	contentTypeJSON := contentTypeSecurity("application/json")

	routes = append(routes,
		post("/product").
			securedWith(contentTypeJSON).
			handle(handler.PostProduct(app)),

		patch("/product/{id:[0-9]+}").
			securedWith(contentTypeJSON).
			handle(handler.PatchProduct(app)),

		delete("/product_variant/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(handler.DeleteProductVariant(app)),

		post("/purchase_order").
			securedWith(contentTypeJSON).
			handle(handler.PostPurchaseOrder(app)),

		delete("/purchase_order/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(handler.DeletePurchaseOrder(app)),

		post("/stock_in").
			securedWith(contentTypeJSON).
			handle(handler.PostStockIn(app)),

		delete("/stock_in/{id:[0-9]+}").
			securedWith(contentTypeJSON).
			handle(handler.DeleteStockIn(app)),

		post("/sales_order").
			securedWith(contentTypeJSON).
			handle(handler.PostSalesOrder(app)),

		delete("/sales_order/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(handler.DeleteSalesOrder(app)),

		post("/stock_out").
			securedWith(contentTypeJSON).
			handle(handler.PostStockOut(app)),

		delete("/stock_out/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(handler.DeleteStockOut(app)),

		// report
		get("/actual_stock").
			securedWith(unsecured()).
			handle(handler.GetActualStock(app)),

		get("/assets_report").
			securedWith(unsecured()).
			handle(handler.GetAssetsReport(app)),

		get("/sales_report").
			securedWith(unsecured()).
			handle(handler.GetSalesReport(app)),

		// csv export
		get("/export/product").
			securedWith(unsecured()).
			handle(handler.GetExportProduct(app)),

		get("/export/stock_in").
			securedWith(unsecured()).
			handle(handler.GetExportStockIn(app)),
	)

	return routes
}
