package router

import (
	"MoP/src/router/route"

	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	// carrega a struct de rotas do pacote route
	return route.RouterConfig(r)
}
