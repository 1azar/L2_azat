package deleteEvent

import (
	"L2-task11/internal/lib/api"
	"log/slog"
	"net/http"
)

type EventDeleter interface {
	DeleteEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventSaver EventDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.deleteEvent.New"
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
