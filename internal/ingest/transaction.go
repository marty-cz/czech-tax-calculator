package ingest

import (
	"fmt"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type TransactionLog struct {
	Purchases         TransactionLogItems
	Sales             TransactionLogItems
	Dividends         TransactionLogItems
	AdditionalIncomes TransactionLogItems
	AdditionalFees    TransactionLogItems
}

type TransactionType int64

const (
	BUY TransactionType = iota
	SELL
	DIVIDEND
	ADDITIONAL_INCOME
	ADDITIONAL_FEE
)

type TransactionLogItem struct {
	// Name of item
	Name string
	// Date of execution
	Date time.Time
	// price per single item
	ItemPrice float64
	// adjusted (due to tax limitation) amount of money send/received from/to a personal bank account
	BankAmount float64
	// non-adjusted amount of money send/received from/to a personal bank account
	OriginalBankAmount float64
	// amount of money used to buy/sell actual item at broker
	BrokerAmount float64
	Fee          float64
	// count of items (event fractions)
	Quantity float64
	// name of Broker who backed the operation
	Broker string
	// Currency used to buy the item (USD, EUR, CZK, ...)
	Currency *util.Currency
	// Exchange rate to CZK in the day of a transaction
	DayExchangeRate float64
	// Exchange rate to CZK in the year of a transaction
	YearExchangeRate float64
	// type of transaction
	Operation TransactionType
	// origin/target country where item was received/buyed
	Country string
}

type TransactionLogItems []*TransactionLogItem

func (items *TransactionLogItems) String() string {
	s := "["
	for i, item := range *items {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%+v", item)
	}
	return s + "]"
}

type newTransactionItem func([]string) (*TransactionLogItem, error)
