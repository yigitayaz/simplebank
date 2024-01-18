package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if c, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(c)
	}
	return false
}
