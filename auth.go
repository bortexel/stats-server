package main

import (
	"context"
	"fmt"
	"net/http"
)

func Authorization(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			header := request.Header.Get("Authorization")
			if header == fmt.Sprintf("Key %s", key) {
				request = request.WithContext(context.WithValue(request.Context(), "authorized", true))
			}

			next.ServeHTTP(writer, request)
		})
	}
}
