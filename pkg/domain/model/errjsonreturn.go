package model

type ErrJsonReturn struct {
	Code   int
	Errors []ErrJson
}

type ErrJson struct {
	Key   string
	Value string
}
