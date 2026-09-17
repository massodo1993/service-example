package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/massodo1993/service-example/order/internal/api/health"
	"github.com/massodo1993/service-example/order/internal/config"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

const handlerTimeout = 30 * time.Second

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(handlerTimeout))
	router.Get("/health", health.Handler)
	router.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP order server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
