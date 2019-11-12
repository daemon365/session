package errors

import "errors"

var (
	ErrSessionNotExist = errors.New("session not exists")
	ErrKeyNotInSession = errors.New("key not in session")
)

