package tax

import (
	"fmt"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type Report struct {
	SellOperations            SellOperations
	TimeTestedItemRevenue     *AccountingValue
	TotalItemRevenue          *AccountingValue
	DividendRevenue           *ValueAndFee
	AdditionalRevenue         *ValueAndFee
	TimeTestedItemFifoExpense *ValueAndFee
	TotalItemFifoExpense      *ValueAndFee
	Year                      time.Time
	Currency                  *util.Currency
}

func (x *Report) String() string {
	return fmt.Sprintf("year: %d sellOpsCount:%d stock:(total:(revenue:(%v) fifoExpense:(%v)) timeTested:(revenue:(%v) fifoExpense:(%v))) dividend:revenue:(%v) additional:revenue:(%v)",
		x.Year.Year(), len(x.SellOperations),
		x.TotalItemRevenue, x.TotalItemFifoExpense,
		x.TimeTestedItemRevenue, x.TimeTestedItemFifoExpense,
		x.DividendRevenue, x.AdditionalRevenue)
}

// map of reports (value) in years (key)
type Reports []*Report
