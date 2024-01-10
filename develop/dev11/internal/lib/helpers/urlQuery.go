package helpers

import (
	"L2-task11/internal/domain"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func ExtractFromPost(r *http.Request) (*domain.ReqParameters, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Распаковка JSON в структуру
	var parms domain.ReqParameters
	if err = json.Unmarshal(body, &parms); err != nil {
		return nil, err
	}

	return &parms, nil
}

func ExtractFromUrl(url *url.URL) (*domain.ReqParameters, error) {
	queryParams := url.Query()

	userId, err := strconv.Atoi(queryParams.Get("user_id"))
	if err != nil {
		return nil, err
	}

	date := queryParams.Get("date")

	return &domain.ReqParameters{UserId: userId, Date: date}, nil
}
