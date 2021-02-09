package routing

import (
	"github.com/chen-keinan/npm-dep-tree/internal/common"
	"github.com/chen-keinan/npm-dep-tree/internal/handler"
	"github.com/chen-keinan/npm-dep-tree/internal/routes"
	"net/http"
	"path"
)

//DependenciesRoutes includes all types required to serve the public routes
type DependenciesRoutes struct {
	DependenciesHandler *handler.Dependencies
	SystemHandler       *handler.System
}

//DepRoutes returns route api for resolving npm dependencies
func (r DependenciesRoutes) DepRoutes() routes.Routes {
	return routes.Routes{
		{
			Name:        "Resolve Npm Dependencies ",
			Method:      http.MethodGet,
			Pattern:     path.Join(common.API, "package-dependencies/{name}/{version}"),
			HandlerFunc: r.DependenciesHandler.ResolveDependencies,
		},
		{
			Name:        "ping check",
			Method:      http.MethodGet,
			Pattern:     path.Join(common.API, "ping"),
			HandlerFunc: r.SystemHandler.Ping,
		},
	}
}
