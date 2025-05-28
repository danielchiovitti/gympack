package helpers

func ToInterfacePtr[T any](v T) *interface{} {
	i := interface{}(v)
	return &i
}
