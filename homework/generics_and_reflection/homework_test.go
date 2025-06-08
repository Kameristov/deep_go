package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	t := reflect.TypeOf(person)
	v := reflect.ValueOf(person)

	var result strings.Builder

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tagStr := field.Tag.Get("properties")
		if tagStr == "" {
			continue
		}

		// Получаем имя поля и опции
		tag, opt, _ := strings.Cut(tagStr, ",")
		if tag == "" {
			tag = field.Name
		}

		// Проверяем наличие опции omitempty
		hasOmitEmpty := false
		for _, opt := range strings.Split(opt, ",") {
			if opt == "omitempty" {
				hasOmitEmpty = true
			}
		}

		// Если поле пустое и есть опция omitempty, то пропускаем его
		if hasOmitEmpty && isEmptyValue(value) {
			continue
		}

		// Форматируем значение в зависимости от типа
		var strValue string

		switch value.Kind() {
		case reflect.String:
			strValue = value.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValue = strconv.Itoa(int(value.Int()))
		case reflect.Bool:
			strValue = strconv.FormatBool(value.Bool())
		}

		result.WriteString(tag)
		result.WriteString("=")
		result.WriteString(strValue)
		result.WriteString("\n")
	}

	return strings.TrimSpace(result.String())
}

// Проверяет, является ли значение пустым
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}
	return false
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
