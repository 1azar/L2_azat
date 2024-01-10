package createEvent

import (
	"L2-task11/internal/http-server/handlers"
	"L2-task11/internal/lib/api"
	"L2-task11/internal/lib/helpers"
	"L2-task11/internal/lib/validators"
	"L2-task11/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=EventSaver
type EventSaver interface {
	SaveEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventSaver EventSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.createEvent.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodPost {
			api.WrapErrorResponse(http.StatusMethodNotAllowed, "only POST allowed", w, lg)
			return
		}

		// парсинг данных из тела post
		eventData, err := helpers.ExtractFromPost(r)
		if err != nil {
			handlers.LogAndErrResponse(http.StatusBadRequest, "invalid data in body", w, lg, err)
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

		// сохранение события
		err = eventSaver.SaveEvent(eventData.UserId, eventData.Date)
		if errors.Is(storage.ErrEventExists, err) {
			handlers.LogAndErrResponse(http.StatusBadRequest, "event already exist", w, lg, err)
			return
		}
		if err != nil {
			handlers.LogAndErrResponse(http.StatusServiceUnavailable, "could not save event", w, lg, err)
			return
		}

		// OK
		api.WrapOkResponse(w, lg)
	}
}
