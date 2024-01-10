package handlers

import (
	"L2-task11/internal/lib/api"
	"log/slog"
	"net/http"
)

func LogAndErrResponse(status int, msg string, w http.ResponseWriter, lg *slog.Logger, externalErr error) {
	lg.Debug("external error", slog.String("error", externalErr.Error()))
	api.WrapErrorResponse(status, msg, w, lg)
}
