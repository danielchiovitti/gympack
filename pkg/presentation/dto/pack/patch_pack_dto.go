package pack

type PatchPackDto struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=150"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
	MaxSize     *int    `json:"maxSize" validate:"omitempty,min=1,max=900000"`
}
