package cmd

import (
	"context"
	"fmt"
	configs "github.com/chen-keinan/npm-dep-tree/config"
	"github.com/chen-keinan/npm-dep-tree/internal/cache"
	"github.com/chen-keinan/npm-dep-tree/internal/handler"
	"github.com/chen-keinan/npm-dep-tree/internal/logger"
	"github.com/chen-keinan/npm-dep-tree/internal/nhttp"
	"github.com/chen-keinan/npm-dep-tree/internal/router"
	"github.com/chen-keinan/npm-dep-tree/internal/router/middleware"
	"github.com/chen-keinan/npm-dep-tree/internal/service"
	"github.com/chen-keinan/npm-dep-tree/internal/workers"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/fx"
	"net"

	"go.uber.org/zap"
	"net/http"
	"time"
)

//RunNpmDepService invoke service
func RunNpmDepService() {
	// start graphite monitor
	InitMonitor()
	// create dependencies injection container
	app := fx.New(
		fx.Provide(logger.NewZapLogger),
		fx.Provide(configs.InitConfig),
		fx.Provide(cache.NewLru),
		fx.Provide(nhttp.NewNpmHTTPClient),
		fx.Provide(service.NewDependencies),
		fx.Provide(handler.NewDependenciesHandler),
		fx.Provide(handler.NewSystemHandler),
		fx.Provide(middleware.NewRateLimitChan),
		fx.Provide(router.NewMuxRouter().RegisterRoutes),
		fx.Invoke(workers.NewWorkerDispatcher().InvokeProcessingWorkers),
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

//InitMonitor start health and metrics service
func InitMonitor() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	if err != nil {
		panic("failed to start monitor")
	}
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)
}
