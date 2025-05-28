package route

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/domain/usecase/pack/create"
	"gympack/pkg/presentation/dto/pack"
	middlewares "gympack/pkg/presentation/middleware"
	errors "gympack/pkg/shared/error"
	"io"
	"net/http"
	"sync"
)

var packRouteLock sync.Mutex
var packRouteInstance *PackRoute

func NewPackRoute(
	dtoValidationMiddleware *middlewares.DtoValidationMiddleware,
	createPackUseCase create.CreatePackUseCaseInterface,
) *PackRoute {
	if packRouteInstance == nil {
		packRouteLock.Lock()
		defer packRouteLock.Unlock()
		if packRouteInstance == nil {
			packRouteInstance = &PackRoute{
				dtoValidationMiddleware: dtoValidationMiddleware,
				createPackUseCase:       createPackUseCase,
			}
		}
	}
	return packRouteInstance
}

type PackRoute struct {
	dtoValidationMiddleware *middlewares.DtoValidationMiddleware
	createPackUseCase       create.CreatePackUseCaseInterface
}

func (p PackRoute) PackRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.With(p.dtoValidationMiddleware.Validate(pack.PostPackDto{})).Post("/", p.postPack)
	return r
}

func (p PackRoute) postPack(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("auth-route")
	ctx, span := tracer.Start(r.Context(), "auth-route.postSignUp")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var dto pack.PostPackDto
	err = json.Unmarshal(body, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message model.PackModel
	err = copier.Copy(&message, &dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := p.createPackUseCase.Execute(ctx, message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	jResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jResp)
}
