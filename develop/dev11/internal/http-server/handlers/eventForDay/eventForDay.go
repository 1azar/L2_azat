package eventForDay

import (
	"L2-task11/internal/lib/api"
	"log/slog"
	"net/http"
)

type DayEventsProvider interface {
	DayEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventSaver DayEventsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.eventForDay.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodGet {
			if err := api.WriteJSON(w, http.StatusMethodNotAllowed, struct{}{}); err != nil {
				lg.Error("could not response to a client")
			}
		}

		lg.Info("implement me")

	}
}
