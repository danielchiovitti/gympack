package shared

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

var lockLogger sync.Mutex
var loggerInstance LoggerInterface

func NewLogger(
	config ConfigInterface,
) LoggerInterface {
	if loggerInstance == nil {
		lockLogger.Lock()
		defer lockLogger.Unlock()
		if loggerInstance == nil {

			//logStashConn, err := net.Dial("tcp", config.GetLogStashUrl())
			//if err != nil {
			//	log.Fatal(err)
			//}

			logrusInstance := log.New()
			//hook := logrustash.New(logStashConn, logrustash.DefaultFormatter(log.Fields{"type": "app_log"}))
			//logrusInstance.Hooks.Add(hook)

			logrusInstance.SetLevel(log.TraceLevel)
			logrusInstance.SetReportCaller(false)
			formatter := &log.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
			}

			logrusInstance.Out = os.Stdout
			logrusInstance.SetFormatter(formatter)

			loggerInstance = &Logger{
				logInstance: logrusInstance,
				config:      config,
			}
		}
	}
	return loggerInstance
}

type Logger struct {
	logInstance *log.Logger
	config      ConfigInterface
}

func (l *Logger) Debug(message string) {
	l.logInstance.Debug(message)
}

func (l *Logger) DebugFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Debug(message)
}

func (l *Logger) Info(message string) {
	l.logInstance.Info(message)
}

func (l *Logger) InfoFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Info(message)
}

func (l *Logger) Warn(message string) {
	l.logInstance.Warn(message)
}

func (l *Logger) WarnFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Warn(message)
}

func (l *Logger) Error(message string) {
	l.logInstance.Error(message)
}

func (l *Logger) ErrorFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Error(message)
}

func (l *Logger) Fatal(message string) {
	l.logInstance.Fatal(message)
}

func (l *Logger) FatalFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Fatal(message)
}

func (l *Logger) Panic(message string) {
	l.logInstance.Panic(message)
}

func (l *Logger) PanicFields(fields map[string]interface{}, message string) {
	l.logInstance.WithFields(fields).Panic(message)
}
