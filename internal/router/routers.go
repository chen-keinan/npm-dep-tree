package router

import (
	"github.com/chen-keinan/npm-dep-tree/internal/handler"
	"github.com/chen-keinan/npm-dep-tree/internal/router/middleware"
	"github.com/chen-keinan/npm-dep-tree/internal/routes"
	"github.com/chen-keinan/npm-dep-tree/internal/routing"
	"github.com/chen-keinan/npm-dep-tree/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
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
func (r MuxRouter) RegisterRoutes(zlog *zap.Logger, depService service.Dep, pr chan middleware.RequestProcessor) *mux.Router {
	dependenciesRouter := &routing.DependenciesRoutes{DependenciesHandler: handler.NewDependenciesHandler(zlog, depService),
		SystemHandler: handler.NewSystemHandler(zlog)}
	r.handlerBuilder(r.router, dependenciesRouter.DepRoutes(), pr)
	return r.router
}

//handlerBuilder build api routes
func (r MuxRouter) handlerBuilder(router *mux.Router, dataIntegrityRoute routes.Routes, pr chan middleware.RequestProcessor) {
	var allRoutes []routes.Routes
	allRoutes = append(allRoutes, dataIntegrityRoute)
	for _, api := range allRoutes {
		for _, route := range api {
			var handler http.Handler
			handler = route.HandlerFunc
			handler = middleware.RequestLimitMiddleware(route.HandlerFunc, pr)
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
}
