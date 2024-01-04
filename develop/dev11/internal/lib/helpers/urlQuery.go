package helpers

import (
	"L2-task11/internal/domain"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func ExtractFromPost(r *http.Request) (*domain.ReqParameters, error) {
	//err := r.ParseForm()
	//if err != nil {
	//	return nil, err
	//}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Закрываем тело запроса после его использования
	defer r.Body.Close()

	// Распаковка JSON в структуру
	var parms domain.ReqParameters
	if err = json.Unmarshal(body, &parms); err != nil {
		return nil, err
	}

	//userId, err := strconv.Atoi(r.FormValue("user_id"))
	//if err != nil {
	//	return nil, err
	//}
	//
	//date, err := time.Parse(time.DateOnly, r.FormValue("date"))
	//if err != nil {
	//	return nil, err
	//}

	return &parms, nil
	//return &domain.ReqParameters{UserId: userId, Date: date}, nil
}

func ExtractFromUrl(url *url.URL) (*domain.ReqParameters, error) {
	//queryParams := url.Query()
	//
	//userId, err := strconv.Atoi(queryParams.Get("user_id"))
	//if err != nil {
	//	return nil, err
	//}
	//
	//date, err := time.Parse(time.DateOnly, queryParams.Get("date"))
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &domain.ReqParameters{UserId: userId, Date: date}, nil
	return nil, nil
}
