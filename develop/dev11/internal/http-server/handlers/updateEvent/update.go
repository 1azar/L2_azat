package updateEvent

import (
	"L2-task11/internal/lib/api"
	"log/slog"
	"net/http"
)

type EventUpdater interface {
	UpdateEvent(userId int, string2 string) error
}

func New(lg *slog.Logger, eventSaver EventUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.updateEvent.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodPost {
			if err := api.WriteJSON(w, http.StatusMethodNotAllowed, struct{}{}); err != nil {
				lg.Error("could not response to a client")
			}
		}

		lg.Info("implement me")

	}
}
