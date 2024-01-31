package appvalidator

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validatorSyncOnce sync.Once
	validate          *validator.Validate
)

func Init() {
	validatorSyncOnce.Do(func() {
		validate = validator.New()
		validate.RegisterTagNameFunc(setTagName)
		validate.RegisterValidation("tmpassword", setPasswordRule)
	})
}

func ValidateStruct(ctx context.Context, object interface{}) (err error) {
	err = validate.StructCtx(ctx, object)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	// errorField := errorFields[0]
	bunchOfErrorMessages := make([]string, len(errorFields))
	for key, field := range errorFields {
		message := fmt.Sprintf("invalid '%s' with value '%v'", field.Field(), field.Value())
		bunchOfErrorMessages[key] = message
	}
	err = fmt.Errorf(strings.Join(bunchOfErrorMessages, ", "))

	return
}

// SetTagName tag name for validator
func setTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func setPasswordRule(fl validator.FieldLevel) bool {
	var (
		hasMinLen  bool   = false
		hasUpper   bool   = false
		hasLower   bool   = false
		hasNumber  bool   = false
		hasSpecial bool   = false
		s          string = fl.Field().String()
	)
	if len(fl.Field().String()) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
