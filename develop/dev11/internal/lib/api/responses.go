package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Result string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ResponseErr2(lg *slog.Logger, w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadGateway)
	lg.Error(msg)
	mb, err := json.Marshal(ErrResponse{Error: msg})
	if err != nil {
		lg.Error("could not marshal response to a client")
	}
	_, err = w.Write(mb)
	if err != nil {
		lg.Error("could not respond to the client")
	}

}
