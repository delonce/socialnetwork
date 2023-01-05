package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	*log.Entry
}
type writerHook struct {
	Writer    []io.Writer
	LogLevels []log.Level
}

var entry *log.Entry

func GetLogger() *Logger {
	return &Logger{entry}
}

func (hook *writerHook) Fire(entry *log.Entry) error {
	logLine, err := entry.String()

	if err != nil {
		return err
	}

	for _, writer := range hook.Writer {
		writer.Write([]byte(logLine))
	}

	return err
}

func (hook *writerHook) Levels() []log.Level {
	return hook.LogLevels
}

func init() {
	logger := log.New()
	logger.SetReportCaller(true)

	logger.Formatter = &log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},

		DisableColors: false,
		FullTimestamp: true,
	}

	logger.SetOutput(io.Discard)
	setHooks(logger)
	logger.SetLevel(log.TraceLevel)

	entry = log.NewEntry(logger)
}

func setHooks(logger *log.Logger) {
	err := os.MkdirAll("logs", 0644)

	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)

	if err != nil {
		panic(err)
	}

	logger.AddHook(&writerHook{
		Writer:    []io.Writer{logFile, os.Stdout},
		LogLevels: log.AllLevels,
	})
}
