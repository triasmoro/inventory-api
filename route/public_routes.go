package route

import (
	"github.com/triasmoro/inventory-api/app"
	"github.com/triasmoro/inventory-api/endpoint"
)

// publicRoutes to register all public endpoints
func publicRoutes(app *app.App) []*HandledRoute {
	var routes []*HandledRoute
	contentTypeJSON := contentTypeSecurity("application/json")

	routes = append(routes,
		post("/product").
			securedWith(contentTypeJSON).
			handle(endpoint.PostProduct(app)),

		patch("/product/{id:[0-9]+}").
			securedWith(contentTypeJSON).
			handle(endpoint.PatchProduct(app)),

		delete("/product_variant/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(endpoint.DeleteProductVariant(app)),

		post("/purchase_order").
			securedWith(contentTypeJSON).
			handle(endpoint.PostPurchaseOrder(app)),

		delete("/purchase_order/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(endpoint.DeletePurchaseOrder(app)),

		post("/stock_in").
			securedWith(contentTypeJSON).
			handle(endpoint.PostStockIn(app)),

		delete("/stock_in/{id:[0-9]+}").
			securedWith(contentTypeJSON).
			handle(endpoint.DeleteStockIn(app)),

		post("/sales_order").
			securedWith(contentTypeJSON).
			handle(endpoint.PostSalesOrder(app)),

		delete("/sales_order/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(endpoint.DeleteSalesOrder(app)),

		post("/stock_out").
			securedWith(contentTypeJSON).
			handle(endpoint.PostStockOut(app)),

		delete("/stock_out/{id:[0-9]+}").
			securedWith(unsecured()).
			handle(endpoint.DeleteStockOut(app)),

		// report
		get("/actual_stock").
			securedWith(unsecured()).
			handle(endpoint.GetActualStock(app)),

		get("/assets_report").
			securedWith(unsecured()).
			handle(endpoint.GetAssetsReport(app)),

		get("/sales_report").
			securedWith(unsecured()).
			handle(endpoint.GetSalesReport(app)),

		// csv export
		get("/export/product").
			securedWith(unsecured()).
			handle(endpoint.GetExportProduct(app)),

		get("/export/stock_in").
			securedWith(unsecured()).
			handle(endpoint.GetExportStockIn(app)),
	)

	return routes
}
