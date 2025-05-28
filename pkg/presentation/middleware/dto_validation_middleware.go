package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gympack/pkg/shared/constant"
	errors2 "gympack/pkg/shared/error"
	"io"
	"net/http"
	"reflect"
)

type DtoValidationMiddleware struct {
}

func NewDtoValidationMiddleware() *DtoValidationMiddleware {
	return &DtoValidationMiddleware{}
}

func (d *DtoValidationMiddleware) Validate(s interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			bodyBytes, err := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			instance := reflect.New(reflect.TypeOf(s)).Interface()
			err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(instance)
			if err != nil {
				http.Error(w, "decoding error JSON", http.StatusBadRequest)
				return
			}

			var validate *validator.Validate
			validate = validator.New(validator.WithRequiredStructEnabled())

			err = validate.Struct(instance)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errors := make(map[string]string)
				for _, err := range err.(validator.ValidationErrors) {
					errors[err.Field()] = err.Tag()
				}
				jsonResponse, err := json.Marshal(errors2.NewDtoValidationErr(http.StatusBadRequest, errors).ErrorJson())
				_, err = w.Write(jsonResponse)
				if err != nil {
					http.Error(w, string(constant.JSON_DECODING_ERROR), http.StatusBadRequest)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
