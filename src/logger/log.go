package logger

import "github.com/borderstech/logmatic"

var logger *logmatic.Logger

func InitialiseLogger() {
	logger = logmatic.NewLogger()
	logger.SetLevel(logmatic.DEBUG)
}

func Info(a string) {
	logger.Info("%s", a)
}

func Warning(a string) {
	logger.Warn("%s", a)
}

func Error(a string) {
	logger.Error("%s", a)
}

func Fatal(a string) {
	logger.Fatal("%s", a)
}

func Debug(a string) {
	logger.Debug("%s", a)
}
