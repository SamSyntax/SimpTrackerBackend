package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request)

func MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := GetApiKey(r.Header)
		if err != nil {
			RespondWithError(w, 403, fmt.Sprintf("auth error: %v", err))
			return
		}

		if apiKey != os.Getenv("API_KEY") {
			RespondWithError(w, 400, fmt.Sprintf("Incorrect key %v", err))
			return
		}
		handler(w, r)
	}
}

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("you must be authenticated to access this endpoint")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
