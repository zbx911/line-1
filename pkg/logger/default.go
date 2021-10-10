package logger

import "github.com/phuslu/log"

func New() *log.Logger {
	logger := &log.Logger{
		Level:      log.TraceLevel,
		TimeFormat: "01-02 15:04:05",
		Writer: &log.MultiWriter{
			InfoWriter:    &log.FileWriter{Filename: "logs/line-api.log", MaxSize: 100 << 20, LocalTime: false},
			ConsoleWriter: &log.ConsoleWriter{ColorOutput: true},
			ConsoleLevel:  log.InfoLevel,
		},
	}
	return logger
}
