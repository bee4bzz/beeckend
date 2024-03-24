package utils

func Or[T any](value, value2 *T) *T {
	if value == nil {
		return value2
	}
	return value
}

func Float64Or(value, value2 *float64) *float64 {
	if value == nil {
		return value2
	}
	return value
}

func StringOr(value, value2 string) string {
	if value == "" {
		return value2
	}
	return value
}

func UintOr(value, value2 uint) uint {
	if value == 0 {
		return value2
	}
	return value
}
