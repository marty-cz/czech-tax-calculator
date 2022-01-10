package ingest

import (
	"fmt"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
	excel "github.com/xuri/excelize/v2"
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
	// amount of money send/received from/to a personal bank account
	BankAmount float64
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

func processSheet(excelFile *excel.File, sheetName string, legend map[string]int, newItemFunction newTransactionItem) (transactions TransactionLogItems, err error) {
	rows, err := excelFile.GetRows(sheetName, excel.Options{RawCellValue: true})
	if err != nil {
		return nil, fmt.Errorf("sheet '%s': %s", sheetName, err)
	}

	transactions = make(TransactionLogItems, 0)
	for rowNo, row := range rows {
		if rowNo == 0 {
			if err := util.ValidateTableHeader(row, legend); err != nil {
				return nil, fmt.Errorf("sheet '%s' at row '%d': %s", sheetName, rowNo, err)
			}
			continue
		} else if util.IsRowEmpty(row, 0) {
			log.Debugf("Sheet '%s' at row '%d' - recognized as empty, skipping", sheetName, rowNo)
			continue
		}

		item, err := newItemFunction(row)
		if err != nil {
			return nil, fmt.Errorf("sheet '%s' at row '%d': %s", sheetName, rowNo, err)
		}

		transactions = append(transactions, item)
	}
	if len(transactions) == 0 {
		log.Warnf("Sheet '%s' has not data to process", sheetName)
	}
	return transactions, nil
}
