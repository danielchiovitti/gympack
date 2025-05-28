package constant

type ErrorType string

const (
	JSON_DECODING_ERROR   ErrorType = "JSON_DECODING_ERROR"
	INVALID_TEMPLATE_DATA ErrorType = "INVALID_TEMPLATE_DATA"
)
