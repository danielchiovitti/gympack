package pack

type PostPackDto struct {
	Name        string `json:"name" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"required,min=3,max=1000"`
	MinSize     uint   `json:"minSize" validate:"required,min=1,max=900000"`
	MaxSize     uint   `json:"maxSize" validate:"required,min=1,max=900000"`
}
