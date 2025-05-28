package presentation

import (
	"github.com/go-chi/chi/v5"
	"gympack/pkg/presentation/route"
	"gympack/pkg/shared"
)

func NewLoader(
	config shared.ConfigInterface,
	packRoute *route.PackRoute,
	logger shared.LoggerInterface,
) *Loader {
	return &Loader{
		Config:    config,
		PackRoute: packRoute,
		Logger:    logger,
	}
}

type Loader struct {
	Config    shared.ConfigInterface
	PackRoute *route.PackRoute
	Logger    shared.LoggerInterface
}

func (l *Loader) GetConfig() shared.ConfigInterface {
	return l.Config
}

func (l *Loader) GetLogger() shared.LoggerInterface {
	return l.Logger
}

func (l *Loader) GetRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/pack", l.PackRoute.PackRoutes())
	//r.Get("/metrics", promhttp.Handler().ServeHTTP)
	return r
}
