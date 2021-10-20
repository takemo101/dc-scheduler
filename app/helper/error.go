package helper

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ErrorsToMap(errors error) map[string]string {
	result := make(map[string]string)
	for key, err := range errors.(validation.Errors) {
		result[key] = err.Error()
	}
	return result
}

func ErrorToMap(key string, err error) map[string]string {
	errors := validation.Errors{}
	errors[key] = err
	return ErrorsToMap(errors)
}
