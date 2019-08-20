package route

// post creates a new POST route. A security handler must be registered next.
func post(tpl string) *Route {
	return &Route{Verb: "POST", Tpl: tpl}
}

// get creates a new GET route. A security handler must be registered next.
func get(tpl string) *Route {
	return &Route{Verb: "GET", Tpl: tpl}
}

// delete creates a new DELETE route. A security handler must be registered next.
func delete(tpl string) *Route {
	return &Route{Verb: "DELETE", Tpl: tpl}
}

// patch creates a new PATCH route. A security handler must be registered next.
func patch(tpl string) *Route {
	return &Route{Verb: "PATCH", Tpl: tpl}
}

// put creates a new PUT route. A security handler must be registered next.
func put(tpl string) *Route {
	return &Route{Verb: "PUT", Tpl: tpl}
}
