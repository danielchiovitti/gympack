package pack

type PostCalcPackDto struct {
	Quantity uint `json:"quantity" validate:"required,min=1,max=900000"`
}
