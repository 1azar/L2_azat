package eventForWeek

import (
	"L2-task11/internal/lib/api"
	"log/slog"
	"net/http"
)

type WeekEventsProvider interface {
	WeekEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventSaver WeekEventsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.eventForWeek.New"
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
