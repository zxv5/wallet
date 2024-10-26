package core

import (
	"database/sql/driver"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateValuer .
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

// GetValidationError .
func GetValidationError(err error, form interface{}) string {
	var vErrors = make(map[string]string)

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ""
	}

	for _, verr := range verrs {
		vErrors[verr.Field()] = verr.Tag()
	}

	t := reflect.TypeOf(form)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var message string
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			sField := t.Field(i)
			msg := sField.Tag.Get("msg")

			if tag, exist := vErrors[sField.Name]; exist {
				message = getTagMessage(msg, tag)
				break
			}
		}
	}

	return message
}

// getTagMessage .
func getTagMessage(msg string, key string) string {
	mapping := make(map[string]string)
	strArr := strings.Split(msg, ",")
	if key == "" {
		key = "__common"
	}

	for _, val := range strArr {
		kv := strings.Split(val, ":")
		if len(kv) == 0 {
			return ""
		}

		if len(kv) >= 2 {
			mapping[kv[0]] = kv[1]
		} else {
			mapping["__common"] = kv[0]
		}
	}

	val, ok := mapping[key]
	if ok {
		return val
	}

	return mapping["__common"]
}

var defaultFuncs = make(map[string]validator.Func)

// RegisterValidation .
func RegisterValidation(validate *validator.Validate) {
	for k, f := range defaultFuncs {
		if err := validate.RegisterValidation(k, f); err != nil {
			panic(err)
		}
	}

	// v.RegisterAlias("id", "email")
}
