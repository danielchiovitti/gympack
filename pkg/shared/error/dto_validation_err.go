package errors

import (
	"errors"
	"fmt"
	"gympack/pkg/domain/model"
	"reflect"
	"strings"
)

type DtoValidationErr struct {
	ErrorCode     int
	ErrorMessages map[string]string
}

func NewDtoValidationErr(errorCode int, errorMessages map[string]string) *DtoValidationErr {
	return &DtoValidationErr{errorCode, errorMessages}
}

func (d DtoValidationErr) Error() string {
	var messages []string
	messages = append(messages, fmt.Sprintf("Code: %d", d.ErrorCode))
	for field, msg := range d.ErrorMessages {
		messages = append(messages, fmt.Sprintf("%s:%s", field, msg))
	}
	return strings.Join(messages, ";")
}

func (d DtoValidationErr) Is(target error) bool {
	var check DtoValidationErr
	if errors.As(target, &check) {
		return reflect.DeepEqual(d, check)
	}
	return false
}

func (d DtoValidationErr) ErrorJson() model.ErrJsonReturn {
	var msg []model.ErrJson
	for key, message := range d.ErrorMessages {
		msg = append(msg, model.ErrJson{Key: key, Value: message})
	}
	return model.ErrJsonReturn{
		Code:   d.ErrorCode,
		Errors: msg,
	}
}
