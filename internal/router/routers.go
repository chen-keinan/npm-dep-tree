package router

import (
	"github.com/chen-keinan/npm-dep-tree/internal/handler"
	"github.com/chen-keinan/npm-dep-tree/internal/router/middleware"
	"github.com/chen-keinan/npm-dep-tree/internal/routes"
	"github.com/chen-keinan/npm-dep-tree/internal/routing"
	"github.com/gorilla/mux"
)

//MuxRouter data
type MuxRouter struct {
	router *mux.Router
}

//NewMuxRouter init new mux router
func NewMuxRouter() MuxRouter {
	muxRouter := MuxRouter{router: mux.NewRouter().StrictSlash(true)}
	return muxRouter
}

//RegisterRoutes instantiate new router
func (r MuxRouter) RegisterRoutes(dependencies *handler.Dependencies, system *handler.System, pr chan middleware.RequestProcessor) *mux.Router {
	dependenciesRouter := &routing.DependenciesRoutes{DependenciesHandler: dependencies, SystemHandler: system}
	r.handlerBuilder(r.router, dependenciesRouter.DepRoutes(), pr)
	return r.router
}

//handlerBuilder build api routes
func (r MuxRouter) handlerBuilder(router *mux.Router, dependenciesRoute routes.Routes, pr chan middleware.RequestProcessor) {
	var allRoutes []routes.Routes
	allRoutes = append(allRoutes, dependenciesRoute)
	for _, api := range allRoutes {
		for _, route := range api {
			handler := middleware.RequestLimitMiddleware(route.HandlerFunc, pr)
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
}
