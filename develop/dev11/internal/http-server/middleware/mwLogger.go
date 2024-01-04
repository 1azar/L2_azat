package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggerWrapper(next http.HandlerFunc, lg *slog.Logger) http.HandlerFunc {
	lg = lg.With(
		slog.String("component", "middleware/logger"))

	fn := func(w http.ResponseWriter, r *http.Request) {
		entry := lg.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
		t1 := time.Now()

		next(w, r)

		entry.Info("request completed",
			slog.String("duration", time.Since(t1).String()),
		)
	}

	return fn
}
