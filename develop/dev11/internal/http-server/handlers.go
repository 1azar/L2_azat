package http_server

import (
	"L2-task11/internal/lib/helpers"
	"net/http"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	params, err := helpers.ExtractFromUrl(r.URL)

}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {

}

func TodayEventHandler(w http.ResponseWriter, r *http.Request) {

}

func WeekEventHandler(w http.ResponseWriter, r *http.Request) {

}

func MonthEventHandler(w http.ResponseWriter, r *http.Request) {

}
