package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ExtractAPIKey extracts an api key from the http headers
// Example:
// Authorization: api_key {value}
func ExtractAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("no authentication provided")
	}

	parts := strings.Split(authorization, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed authorization")
	}

	if parts[0] != "api_key" {
		return "", errors.New("authorization header value should be in the format `Authorization: api_key {value}`")
	}

	return parts[1], nil
}
