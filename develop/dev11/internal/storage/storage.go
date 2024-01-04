package storage

import "errors"

var (
	ErrEventExists    = errors.New("event already exist")
	ErrEventNotExists = errors.New("event does not exist")
)
