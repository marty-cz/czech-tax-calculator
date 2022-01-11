package tax

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type AccountingValue struct {
	ValueWithDayExchangeRate  float64
	ValueWithYearExchangeRate float64
	Currency                  *util.Currency
}

func newAccountingValue(valueWithDayRate, valueWithYearRate float64, c *util.Currency) *AccountingValue {
	if c == nil {
		c = DEFAULT_CURRENCY
	}
	return &AccountingValue{ValueWithDayExchangeRate: valueWithDayRate, ValueWithYearExchangeRate: valueWithYearRate, Currency: c}
}

func (x *AccountingValue) String() string {
	return fmt.Sprintf("withDayExchange:%s %v withYearExchange:%s %v", x.Currency.Symbol, x.ValueWithDayExchangeRate, x.Currency.Symbol, x.ValueWithYearExchangeRate)
}

func (x *AccountingValue) Add(add *AccountingValue) {
	x.ValueWithDayExchangeRate += add.ValueWithDayExchangeRate
	x.ValueWithYearExchangeRate += add.ValueWithYearExchangeRate
}

func (x *AccountingValue) Sub(add *AccountingValue) {
	x.ValueWithDayExchangeRate -= add.ValueWithDayExchangeRate
	x.ValueWithYearExchangeRate -= add.ValueWithYearExchangeRate
}
func (x *AccountingValue) MultiplyNew(multiplicator float64) *AccountingValue {
	return &AccountingValue{
		ValueWithDayExchangeRate:  x.ValueWithDayExchangeRate * multiplicator,
		ValueWithYearExchangeRate: x.ValueWithYearExchangeRate * multiplicator,
		Currency:                  x.Currency,
	}
}
