package logger

import (
	"log/slog"
	"os"
)

var L *slog.Logger

func Init(format string) {
	var handler slog.Handler

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	L = slog.New(handler)
	slog.SetDefault(L)

	L.Info("logger initialized", "format", format)
}
