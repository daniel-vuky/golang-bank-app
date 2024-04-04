package api

import (
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/go-playground/validator/v10"
)

// validCurrency func to valid if input currency is valid
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportCurrency(currency)
	}
	return false
}
