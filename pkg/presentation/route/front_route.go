package route

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"gympack/pkg/domain/usecase/pack/get"
	errors "gympack/pkg/shared/error"
	"gympack/pkg/shared/helpers"
	"html/template"
	"net/http"
)

//go:embed templates/index.html
//go:embed templates/menu.html
//go:embed templates/view_packs.html
//go:embed templates/create_pack.html
var tmplFS embed.FS

func NewFrontRoute(
	getPackUseCase get.GetPackUseCaseInterface,
) *FrontRoute {
	return &FrontRoute{
		getPackUseCase: getPackUseCase,
	}
}

type FrontRoute struct {
	getPackUseCase get.GetPackUseCaseInterface
}

func (f *FrontRoute) FrontRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", f.getMain)
	r.Get("/view-packs", f.viewPacks)
	r.Get("/create-pack", f.createPack)
	return r
}

func (f *FrontRoute) getMain(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(tmplFS,
		"templates/index.html",
		"templates/menu.html",
		"templates/create_pack.html",
	))

	err := tmpl.Execute(w, nil)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
}

func (f *FrontRoute) viewPacks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(tmplFS,
		"templates/view_packs.html",
	))

	packs, err := f.getPackUseCase.Execute(r.Context(), "")

	err = tmpl.Execute(w, packs)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
}

func (f *FrontRoute) createPack(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(tmplFS,
		"templates/create_pack.html",
	))

	err := tmpl.Execute(w, nil)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
}
