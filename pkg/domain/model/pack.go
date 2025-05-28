package model

type PackModel struct {
	Id          string
	Name        string
	Description string
	MinSize     int
	MaxSize     int
	BaseModel
}
