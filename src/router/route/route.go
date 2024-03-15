package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
}

func RouterConfig(r *mux.Router) *mux.Router {
	// carrega a array de rotas do pacote route.agents.go
	routes := agentsRoute

	for _, route := range routes {

		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}
