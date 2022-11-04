package main

import (
	"errors"
	"fmt"
	"net/http"
)

func Authorization(key string) func(handler ActionHandler) ActionHandler {
	return func(next ActionHandler) ActionHandler {
		return func(r *http.Request, body []byte) (any, error, int) {
			header := r.Header.Get("Authorization")
			if header != fmt.Sprintf("Key %s", key) {
				return nil, errors.New("unauthorized"), http.StatusUnauthorized
			}

			return next(r, body)
		}
	}
}
