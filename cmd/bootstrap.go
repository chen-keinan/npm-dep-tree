package cmd

import (
	"context"
	"fmt"
	configs "github.com/chen-keinan/npm-dep-tree/config"
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/handler"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/internal/router"
	"github.com/chen-keinan/npm-dep-tree/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//RunNpmDepService invoke service
func RunNpmDepService() {
	// create dependencies injection container
	app := fx.New(
		fx.Provide(logger.NewZapLogger),
		fx.Provide(configs.InitConfig),
		fx.Provide(cache.NewLru),
		fx.Provide(service.NewDependencies),
		fx.Provide(handler.NewDependenciesHandler),
		fx.Provide(handler.NewSystemHandler),
		fx.Provide(router.NewMuxRouter().RegisterRoutes),
		fx.Invoke(runHTTPServer),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

// start http server
func runHTTPServer(lifecycle fx.Lifecycle, routes *mux.Router, c *configs.Config, zlog *zap.Logger) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error {
		httpPort := ":" + c.GetStringValue("HttpPort")
		zlog.Info(fmt.Sprintf("HTTP port: %s", httpPort))
		srv := &http.Server{
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Minute,
			Addr:         httpPort,
			Handler:      routes}
		return srv.ListenAndServe()
	},
	})
}

/*
//StartSystemServices start health and metrics service
func StartSystemServices() {
	//http.Handle("/_/ping", handlers.HandleHealthCheck())
	//http.Handle("/_/metrics.json", reporters.HandleMetricsJson())
	adminPort := ":" + configs.GetStringValue("AdminPort")
	log.Infof("Admin port: %s", adminPort)
	if err := http.ListenAndServe(adminPort, nil); err != nil {
		panic("failed to start health and metrics services")
	}
}
*/
