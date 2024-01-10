package createEvent_test

import (
	"L2-task11/internal/domain"
	createEvent "L2-task11/internal/http-server/handlers/createEvent"
	"L2-task11/internal/http-server/handlers/createEvent/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEventHandler(t *testing.T) {
	tests := []struct {
		name      string
		userId    int
		date      string
		respErr   string
		mockError error
	}{
		{name: "success", userId: 0, date: "2020-09-11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eventSaverMock := mocks.NewEventSaver(t)

			if tt.respErr == "" || tt.mockError != nil {
				eventSaverMock.On("SaveEvent", tt.userId, tt.date).
					Return(tt.mockError).
					Once()
			}

			lg := slog.New(slog.NewJSONHandler(io.Discard, nil))
			handler := createEvent.New(lg, eventSaverMock)

			input := fmt.Sprintf(`{"user_id": %d, "date":"%s"}`, tt.userId, tt.date)

			req, err := http.NewRequest(http.MethodPost, "/create_event", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp domain.ReqParameters

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			//require.Equal(t, tt.respErr, resp.Error)

			// TODO закончить

		})
	}
}
