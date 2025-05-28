package mappers

type MapperInterface[T, U any] interface {
	ToEntity(model T) (*U, error)
	ToModel(entity U) (*T, error)
}
