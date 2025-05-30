package helpers

import (
	"errors"
	"reflect"
)

var ErrInvalidDestination = errors.New("destination must be a non-nil pointer to a struct")

func PatchStruct(dst interface{}, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return ErrInvalidDestination
	}
	dstVal = dstVal.Elem()

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcVal.Type().Field(i).Name)

		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		if srcField.Kind() == reflect.Ptr {
			if !srcField.IsNil() {
				dstField.Set(srcField.Elem())
			}
		} else {
			zero := reflect.Zero(srcField.Type())
			if !reflect.DeepEqual(srcField.Interface(), zero.Interface()) {
				dstField.Set(srcField)
			}
		}
	}

	return nil
}
