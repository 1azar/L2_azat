package createEvent

import (
	"L2-task11/internal/lib/api"
	"L2-task11/internal/lib/helpers"
	"L2-task11/internal/storage"
	"errors"
	"log/slog"
	"net/http"
)

type EventSaver interface {
	SaveEvent(userId int, eventDate string) error
}

func New(lg *slog.Logger, eventSaver EventSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.createEvent.New"
		lg = lg.With(slog.String("fn", fn))

		// check method
		if r.Method != http.MethodPost {
			if err := api.WriteJSON(w, http.StatusMethodNotAllowed, struct{}{}); err != nil {
				lg.Error("could not response to a client")
			}
		}

		// данные из тела post
		eventData, err := helpers.ExtractFromPost(r)
		if err != nil {
			err = api.WriteJSON(w, http.StatusBadRequest, api.ErrResponse{Error: "invalid form"})
			if err != nil {
				lg.Error("could not response to a client: ", err)
			}
			return
		}
		lg.Debug("request url query decoded")

		// сохранение события
		err = eventSaver.SaveEvent(eventData.UserId, eventData.Date)
		if errors.Is(storage.ErrEventExists, err) {
			err = api.WriteJSON(w, http.StatusBadRequest, api.ErrResponse{Error: "event already exist"})
			if err != nil {
				lg.Error("could not response to a client:", err)
			}
			return
		}
		if err != nil {
			err = api.WriteJSON(w, http.StatusServiceUnavailable, api.ErrResponse{Error: "could not save event"})
			if err != nil {
				lg.Error("could not response to a client:", err)
			}
			return
		}

	}
}
