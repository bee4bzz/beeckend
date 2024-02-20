package utils

import "github.com/labstack/gommon/random"

func String(s string) *string {
	return &s
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
