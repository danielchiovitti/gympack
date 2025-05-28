package shared

type LoggerInterface interface {
	Debug(message string)
	DebugFields(fields map[string]interface{}, message string)
	Info(message string)
	InfoFields(fields map[string]interface{}, message string)
	Warn(message string)
	WarnFields(fields map[string]interface{}, message string)
	Error(message string)
	ErrorFields(fields map[string]interface{}, message string)
	Fatal(message string)
	FatalFields(fields map[string]interface{}, message string)
	Panic(message string)
	PanicFields(fields map[string]interface{}, message string)
}
