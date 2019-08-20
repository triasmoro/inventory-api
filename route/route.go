package route

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/triasmoro/inventory-api/app"
)

// SecurityHandler is a kind of middleware that satisfies security requirements
type SecurityHandler func(http.Handler) http.Handler

// Route is an incomplete Route comprising only verb and path (as a gorilla/mux template). It must
// next be `SecuredWith`.
type Route struct {
	Verb string
	Tpl  string // template
}

// HandledRoute is a fully defined route. It is ready to be `Attach`d.
type HandledRoute struct {
	*SecuredRoute
	handler http.Handler
}

// SecuredRoute is an incomplete Route with a defined SecurityHandler. It is ready for a
// http.Handler.
type SecuredRoute struct {
	Route
	security SecurityHandler
}

// Router handler
func Router(app *app.App) http.Handler {
	r := mux.NewRouter()
	attach(r, publicRoutes(app)...)

	return wrapRouter(r)
}

func (hr *HandledRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	hr.security(hr.handler).ServeHTTP(response, request)
}

// attach is the adapter for adding HandledRoutes to a gorilla/mux Router.
func attach(router *mux.Router, routes ...*HandledRoute) {
	for _, r := range routes {
		router.
			Methods(r.Verb).
			Path(r.Tpl).
			Handler(r)
	}
}

// add middlewares here
func wrapRouter(r *mux.Router) http.Handler {
	stack := handlers.CombinedLoggingHandler(os.Stdout, r)

	return stack
}

// securedWith registers a security handler for a route. A handler must be registered next.
func (r Route) securedWith(fn SecurityHandler) *SecuredRoute {
	return &SecuredRoute{r, fn}
}

// handle registers a HandlerFunc. The route may now be `Attach`d.
func (r *SecuredRoute) handle(h http.Handler) *HandledRoute {
	return &HandledRoute{r, h}
}
