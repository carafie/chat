package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/carafie/chat/api/internal/web"
	"github.com/carafie/chat/api/pkg/pubsub"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindTime {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	}))

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "127.0.0.1:8080"
	}

	pubsubClient := pubsub.New()
	webHandler := web.NewHandler(pubsubClient)
	webRouter := web.NewRouter(webHandler)
	webServer := http.Server{
		Addr:     address,
		Handler:  webRouter,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("start server", slog.String("address", address))
	if err := webServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("stop server")
		} else {
			logger.Error("stop server", slog.String("error", err.Error()))
		}
	}
}
