package tax

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type ValueAndFee struct {
	Value *AccountingValue
	Fee   *AccountingValue
}

func newEmptyValueAndFee(c *util.Currency) *ValueAndFee {
	if c == nil {
		c = DEFAULT_CURRENCY
	}
	return &ValueAndFee{
		Value: &AccountingValue{ValueWithDayExchangeRate: 0, ValueWithYearExchangeRate: 0, Currency: c},
		Fee:   &AccountingValue{ValueWithDayExchangeRate: 0, ValueWithYearExchangeRate: 0, Currency: c},
	}
}

func (x *ValueAndFee) String() string {
	return fmt.Sprintf("val:(%v) fee:(%v)", x.Value, x.Fee)
}
