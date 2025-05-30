package presentation

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"gympack/pkg/presentation/route"
	"gympack/pkg/shared"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed route/static/*
var staticFiles embed.FS

func NewLoader(
	config shared.ConfigInterface,
	packRoute *route.PackRoute,
	frontRoute *route.FrontRoute,
	logger shared.LoggerInterface,
) *Loader {
	return &Loader{
		Config:     config,
		PackRoute:  packRoute,
		Logger:     logger,
		FrontRoute: frontRoute,
	}
}

type Loader struct {
	Config     shared.ConfigInterface
	PackRoute  *route.PackRoute
	Logger     shared.LoggerInterface
	FrontRoute *route.FrontRoute
}

func (l *Loader) GetConfig() shared.ConfigInterface {
	return l.Config
}

func (l *Loader) GetLogger() shared.LoggerInterface {
	return l.Logger
}

func (l *Loader) GetRoutes() *chi.Mux {
	r := chi.NewRouter()

	subStaticFiles, err := fs.Sub(staticFiles, "route/static")
	if err != nil {
		panic(err)
	}
	FileServer(r, "/static", http.FS(subStaticFiles))

	r.Mount("/pack", l.PackRoute.PackRoutes())
	r.Mount("/front", l.FrontRoute.FrontRoutes())
	//r.Get("/metrics", promhttp.Handler().ServeHTTP)
	return r
}
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	r.Get(path+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
