package eventForMonth

import (
	"L2-task11/internal/domain"
	"L2-task11/internal/http-server/handlers"
	"L2-task11/internal/lib/api"
	"L2-task11/internal/lib/helpers"
	"L2-task11/internal/lib/validators"
	"log/slog"
	"net/http"
)

type EventsForMonthProvider interface {
	EventsForMonth(userId int, string2 string) (*[]domain.ReqParameters, error)
}

func New(lg *slog.Logger, eventsForMonth EventsForMonthProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.eventForMonth.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodGet {
			api.WrapErrorResponse(http.StatusMethodNotAllowed, "method not allowed", w, lg)
			return
		}

		// парсинг данных из урла
		eventData, err := helpers.ExtractFromUrl(r.URL)
		if err != nil {
			handlers.LogAndErrResponse(http.StatusBadRequest, "invalid url query", w, lg, err)
			return
		}

		// валидация даты
		if !validators.IsValidDate(eventData.Date) {
			api.WrapErrorResponse(http.StatusBadRequest, "invalid date format", w, lg)
			return
		}

		// валидация user_id
		if eventData.UserId < 0 {
			api.WrapErrorResponse(http.StatusBadRequest, "invalid user_id", w, lg)
			return
		}

		// получение события
		resp, err := eventsForMonth.EventsForMonth(eventData.UserId, eventData.Date)
		if err != nil {
			handlers.LogAndErrResponse(http.StatusServiceUnavailable, "could not gather events", w, lg, err)
			return
		}

		// OK
		if err = api.WriteJSON(w, http.StatusOK, resp); err != nil {
			lg.Error("could not write JSON to client", slog.String("error", err.Error()))
		}

	}
}
