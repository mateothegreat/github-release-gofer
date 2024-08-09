package config

import "errors"

func GetIntFromSlice(min int, max int, values ...int) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("no values provided")
	}

	for _, value := range values {
		if value >= min && value <= max {
			return value, nil
		}
	}

	return 0, errors.New("value not in range")
}

func GetStringFromSlice(values ...string) (string, error) {
	if len(values) == 0 {
		return "", errors.New("no values provided")
	}

	for _, value := range values {
		if value != "" {
			return value, nil
		}
	}

	return "", errors.New("no non-empty values provided")
}
