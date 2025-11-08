package auth

import (
	"errors"
	"net/http"
	"strings"
)

// It extracts API Keys from Header of HTTP Request
// If can find API Key return as string else error is returned.
// As Author of the server as can decide how the Authentication header to looks like
// Example:
// Authorization: ApiKey {insert the Apikey here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("No authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("Malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Malformed first part of the auth header")
	}

	return vals[1], nil

}
