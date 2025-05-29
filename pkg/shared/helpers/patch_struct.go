package helpers

import "reflect"

func PatchStruct(dst interface{}, src interface{}) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcVal.Type().Field(i).Name)

		if srcField.Kind() == reflect.Ptr && !srcField.IsNil() && dstField.CanSet() {
			dstField.Set(srcField.Elem())
		}
	}
}
