package utils

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
