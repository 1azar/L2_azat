package updateEvent

import (
	"L2-task11/internal/http-server/handlers"
	"L2-task11/internal/lib/api"
	"L2-task11/internal/lib/helpers"
	"L2-task11/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type EventUpdater interface {
	UpdateEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventUpdater EventUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.updateEvent.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodPost {
			api.WrapErrorResponse(http.StatusMethodNotAllowed, "only POST allowed", w, lg)
			return
		}

		// extract data from post body
		eventData, err := helpers.ExtractFromPost(r)
		if err != nil {
			handlers.LogAndErrResponse(http.StatusBadRequest, "invalid data in body", w, lg, err)
			return
		}

		// обновление события
		err = eventUpdater.UpdateEvent(eventData.UserId, eventData.Date)
		if errors.Is(storage.ErrEventNotExists, err) {
			handlers.LogAndErrResponse(http.StatusBadRequest, "event not exist", w, lg, err)
			return
		}
		if err != nil {
			handlers.LogAndErrResponse(http.StatusServiceUnavailable, "could not update event", w, lg, err)
			return
		}

		// OK
		api.WrapOkResponse(w, lg)

	}
}
