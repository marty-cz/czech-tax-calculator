package tax

import (
	"fmt"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type Report struct {
	SellOperations        SellOperations
	TimeTestedItemRevenue *AccountingValue
	TotalItemRevenue      *AccountingValue
	// map of dividends (value) in countries (key)
	DividendReports           map[string]*DividendReport
	AdditionalRevenue         *ValueAndFee
	TimeTestedItemFifoExpense *ValueAndFee
	TotalItemFifoExpense      *ValueAndFee
	Year                      time.Time
	Currency                  *util.Currency
}

func (x *Report) String() string {
	return fmt.Sprintf("year: %d sellOpsCount:%d stock:(total:(revenue:(%v) fifoExpense:(%v)) timeTested:(revenue:(%v) fifoExpense:(%v))) dividend:[%v]) additional:revenue:(%v)",
		x.Year.Year(), len(x.SellOperations),
		x.TotalItemRevenue, x.TotalItemFifoExpense,
		x.TimeTestedItemRevenue, x.TimeTestedItemFifoExpense,
		x.DividendReports,
		x.AdditionalRevenue)
}

// map of reports (value) in years (key)
type Reports []*Report

type DividendReport struct {
	RawRevenue *ValueAndFee
	PaidTax    *AccountingValue
	Country    string
}

func (x *DividendReport) String() string {
	return fmt.Sprintf("country:%v rawRevenue:(%v) paidTax:(%v)",
		x.Country, x.RawRevenue, x.PaidTax)
}
