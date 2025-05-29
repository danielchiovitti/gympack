package route

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/domain/usecase/pack/create"
	delete2 "gympack/pkg/domain/usecase/pack/delete"
	"gympack/pkg/domain/usecase/pack/get"
	"gympack/pkg/domain/usecase/pack/update"
	"gympack/pkg/presentation/dto/pack"
	middlewares "gympack/pkg/presentation/middleware"
	errors "gympack/pkg/shared/error"
	"gympack/pkg/shared/helpers"
	"io"
	"net/http"
	"sync"
)

var packRouteLock sync.Mutex
var packRouteInstance *PackRoute

func NewPackRoute(
	dtoValidationMiddleware *middlewares.DtoValidationMiddleware,
	createPackUseCase create.CreatePackUseCaseInterface,
	getPackUseCase get.GetPackUseCaseInterface,
	deletePackUseCase delete2.DeletePackUseCaseInterface,
	updatePackUseCase update.UpdatePackUseCaseInterface,
) *PackRoute {
	if packRouteInstance == nil {
		packRouteLock.Lock()
		defer packRouteLock.Unlock()
		if packRouteInstance == nil {
			packRouteInstance = &PackRoute{
				dtoValidationMiddleware: dtoValidationMiddleware,
				createPackUseCase:       createPackUseCase,
				getPackUseCase:          getPackUseCase,
				deletePackUseCase:       deletePackUseCase,
				updatePackUseCase:       updatePackUseCase,
			}
		}
	}
	return packRouteInstance
}

type PackRoute struct {
	dtoValidationMiddleware *middlewares.DtoValidationMiddleware
	createPackUseCase       create.CreatePackUseCaseInterface
	getPackUseCase          get.GetPackUseCaseInterface
	deletePackUseCase       delete2.DeletePackUseCaseInterface
	updatePackUseCase       update.UpdatePackUseCaseInterface
}

func (p PackRoute) PackRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.With(p.dtoValidationMiddleware.Validate(pack.PostPackDto{})).Post("/", p.postPack)
	r.Get("/", p.getPack)
	r.Get("/{id}", p.getPack)
	r.Delete("/{id}", p.deletePack)
	r.Patch("/{id}", p.patchPack)
	return r
}

func (p PackRoute) postPack(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("pack-route")
	ctx, span := tracer.Start(r.Context(), "pack-route.post")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	defer r.Body.Close()

	var dto pack.PostPackDto
	err = json.Unmarshal(body, &dto)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var pModel model.PackModel
	err = copier.Copy(&pModel, &dto)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := p.createPackUseCase.Execute(ctx, pModel)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	helpers.JsonResponse(w, http.StatusOK, resp)
}

func (p PackRoute) getPack(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("pack-route")
	ctx, span := tracer.Start(r.Context(), "pack-route.get")
	defer span.End()

	id := chi.URLParam(r, "id")
	resp, err := p.getPackUseCase.Execute(ctx, id)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	helpers.JsonResponse(w, http.StatusOK, resp)
}

func (p PackRoute) deletePack(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("pack-route")
	ctx, span := tracer.Start(r.Context(), "pack-route.delete")
	defer span.End()

	id := chi.URLParam(r, "id")
	resp, err := p.deletePackUseCase.Execute(ctx, id)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	helpers.JsonResponse(w, http.StatusOK, resp)
}

func (p PackRoute) patchPack(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("pack-route")
	ctx, span := tracer.Start(r.Context(), "pack-route.delete")
	defer span.End()

	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{Error: "id required"})
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	defer r.Body.Close()

	var dto pack.PatchPackDto
	err = json.Unmarshal(body, &dto)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	pModel, err := p.getPackUseCase.Execute(ctx, id)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(*pModel) == 0 {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: "id not found",
		})
		return
	}

	uModel := (*pModel)[0]
	helpers.PatchStruct(&uModel, &dto)

	resp, err := p.updatePackUseCase.Execute(ctx, id, uModel)
	if err != nil {
		helpers.JsonResponse(w, http.StatusBadRequest, errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	helpers.JsonResponse(w, http.StatusOK, resp)
}
