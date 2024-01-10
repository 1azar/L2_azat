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
	Result string `json:"result"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WrapOkResponse(w http.ResponseWriter, lg *slog.Logger) {
	err := WriteJSON(w, http.StatusOK, "OK")
	if err != nil {
		lg.Error("could not response to a client:", err)
	}
}

func WrapErrorResponse(status int, msg string, w http.ResponseWriter, lg *slog.Logger) {
	err := WriteJSON(w, status, ErrResponse{Error: msg})
	if err != nil {
		lg.Error("could not response to a client:", err)
	}
}
