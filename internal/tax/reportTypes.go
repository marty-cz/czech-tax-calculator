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
	// map of dividends per broker (value) in countries (key)
	DividendReports           map[string]*BrokerDividendReports
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

// map of dividends (value) in brokers (key)
type BrokerDividendReports map[string]*DividendReport

func (m BrokerDividendReports) GetAll() map[string]*DividendReport {
	return m
}

func (m BrokerDividendReports) Get(broker string) (*DividendReport, bool) {
	report, exists := m[broker]
	return report, exists
}

func (m BrokerDividendReports) Set(broker string, report *DividendReport) error {
	m[broker] = report
	return nil
}

type DividendReport struct {
	RawRevenue         *ValueAndFee
	PaidTax            *AccountingValue
	OriginalRawRevenue *ValueAndFee
	OriginalPaidTax    *AccountingValue
	Country    string
	Broker     string
}

func (x *DividendReport) String() string {
	return fmt.Sprintf("country:%v broker:%v rawRevenue:(%v) paidTax:(%v)",
		x.Country, x.Broker, x.RawRevenue, x.PaidTax)
}
