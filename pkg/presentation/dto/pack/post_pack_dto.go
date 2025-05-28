package pack

type PostPackDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	MinSize     uint   `json:"minSize" validate:"min=1,max=900000"`
	MaxSize     uint   `json:"maxSize" validate:"min=1,max=900000"`
}
