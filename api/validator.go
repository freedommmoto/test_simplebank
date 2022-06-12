package api

import (
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(ft validator.FieldLevel) bool {
	if currency, ok := ft.Field().Interface().(string); ok {
		return tool.IsSupportedCurrencyOrNot(currency)
	}
	return false
}
