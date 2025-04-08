package log

import (
	"io"
	"log"
	"log/slog"
	"os"
	"time"
)

type LoggerConfig struct {
	Format    string `env:"LOGGING_FORMAT" yaml:"format" env-default:"text"`
	Level     string `env:"LOGGING_LEVEL" yaml:"level" env-default:"info"`
	Directory string `env:"LOGGING_DIRECTORY" yaml:"directory" env-default:""`
}

func NewLogger(cfg LoggerConfig) (*slog.Logger, *os.File) {

	var level slog.Level
	switch cfg.Level {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	writeToFile := len(cfg.Directory) != 0
	var file *os.File

	var writer io.Writer
	if writeToFile {
		file, err := createLogFile(cfg.Directory)
		if err != nil {
			log.Printf("error while creating log file: %s\n", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
	} else {
		writer = os.Stdout
	}

	var handler slog.Handler
	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	case "json":
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler), file
}

func createLogFile(dirPath string) (*os.File, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, err
		}
	}

	filePath := dirPath + "/" + time.Now().String() + ".log"

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return logFile, nil

}

func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
