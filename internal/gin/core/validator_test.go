package core

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required" msg:"Name is required"`
	Age  int    `validate:"gte=0,lte=130" msg:"Age must be between 0 and 130"`
}

func TestValidateValuer(t *testing.T) {
	field := reflect.ValueOf(TestStruct{Name: "Test"})
	val := ValidateValuer(field)
	assert.Nil(t, val)
}

func TestGetValidationError(t *testing.T) {
	validate := validator.New()
	form := &TestStruct{}

	// Triggering validation error
	err := validate.Struct(form)
	assert.Error(t, err)

	message := GetValidationError(err, form)
	assert.Contains(t, message, "Name is required")

	// Check for age validation
	form.Name = "John"
	form.Age = 150
	err = validate.Struct(form)
	assert.Error(t, err)

	message = GetValidationError(err, form)
	assert.Contains(t, message, "Age must be between 0 and 130")
}

func TestGetTagMessage(t *testing.T) {
	msg := "required:This field is required, email:Invalid email format"
	key := "required"
	expected := "This field is required"

	result := getTagMessage(msg, key)
	assert.Equal(t, expected, result)
}
