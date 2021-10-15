package helper

import (
	"reflect"
)

// StructToMap to map
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}

// StructToJsonMap to map by json tag
func StructToJsonMap(data interface{}) map[string]interface{} {
	return StructToTagMap("json", data)
}

// StructToFormMap to map by form tag
func StructToFormMap(data interface{}) map[string]interface{} {
	return StructToTagMap("form", data)
}

// StructToFormMap to map by string tag
func StructToTagMap(tag string, data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Tag.Get(tag)
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}
