package utils

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/random"
)

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func Float64(f float64) *float64 {
	return &f
}

func ValidEmail() string {
	return random.String(10) + "@example.com"
}

func ValidName() string {
	return random.String(10)
}

func BearerAuthHeader(header *http.Header, token string) {
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}
