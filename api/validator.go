package api

import (
	"db/db/util"

	"github.com/go-playground/validator/v10"
)

//Cách custom lại các trường validate binding

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		return util.IssupportedCurrency(currency)
	}
	return false
}
