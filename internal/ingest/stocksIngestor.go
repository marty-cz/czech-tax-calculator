package ingest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
	excel "github.com/xuri/excelize/v2"
)

var (
	BUY_TABLE_LEGEND = map[string]int{
		"STOCK":       0,
		"DATE":        1,
		"STOCK PRICE": 2,
		"PAID":        3,
		"FEE":         4,
		"AMOUNT":      5,
		"QUANTITY":    6,
		"BROKER":      7,
		"CURRENCY":    8,
	}
	SELL_TABLE_LEGEND = map[string]int{
		"STOCK":       0,
		"DATE":        1,
		"STOCK PRICE": 2,
		"RECEIVED":    3,
		"FEE":         4,
		"AMOUNT":      5,
		"QUANTITY":    6,
		"BROKER":      7,
		"CURRENCY":    8,
	}
	DIVIDEND_TABLE_LEGEND = map[string]int{
		"STOCK":    0,
		"DATE":     1,
		"RECEIVED": 2,
		"AMOUNT":   3,
		"PAID TAX": 4,
		"BROKER":   5,
		"CURRENCY": 6,
	}
)

type TransactionType int64

const (
	BUY TransactionType = iota
	SELL
	DIVIDEND
)

type TransactionLogItem struct {
	// name of item
	name string
	// date of execution
	date time.Time
	// price per single item
	itemPrice float64
	// amount of money send/received from/to a personal bank account
	bankAmount float64
	// amount of money used to buy/sell actual item at broker
	brokerAmount float64
	fee          float64
	// count of items (event fractions)
	quantity float64
	// name of broker who backed the operation
	broker string
	// currency used to buy the item (USD, EUR, CZK, ...)
	currency string
	// type of transaction
	operation TransactionType
}

type newTransactionItem func([]string) (*TransactionLogItem, error)

func newBuyItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		name:      row[BUY_TABLE_LEGEND["STOCK"]],
		broker:    row[BUY_TABLE_LEGEND["BROKER"]],
		operation: BUY,
	}

	if rawDate, err := strconv.ParseFloat(row[BUY_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.itemPrice, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %s", err)
	}
	if item.bankAmount, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["PAID"]], 64); err != nil {
		return nil, fmt.Errorf("paid is not a number: %s", err)
	}
	if item.brokerAmount, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.fee, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.quantity, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.currency, err = util.ConvertCurrency(row[BUY_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	return &item, nil
}

func newSellItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		name:      row[SELL_TABLE_LEGEND["STOCK"]],
		broker:    row[SELL_TABLE_LEGEND["BROKER"]],
		operation: SELL,
	}

	if rawDate, err := strconv.ParseFloat(row[SELL_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.itemPrice, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %s", err)
	}
	if item.bankAmount, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %s", err)
	}
	if item.brokerAmount, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.fee, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.quantity, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.currency, err = util.ConvertCurrency(row[SELL_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	return &item, nil
}

func newDividendItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		name:      row[DIVIDEND_TABLE_LEGEND["STOCK"]],
		broker:    row[DIVIDEND_TABLE_LEGEND["BROKER"]],
		operation: DIVIDEND,
	}

	if rawDate, err := strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.bankAmount, err = strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %s", err)
	}
	if item.brokerAmount, err = strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.currency, err = util.ConvertCurrency(row[DIVIDEND_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	return &item, nil
}

func ProcessStocks(filePath string) (err error) {

	log.Debugf("Processing stocks input file '%s'", filePath)

	f, err := excel.OpenFile(filePath)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Errorf("Cannot close file '%s' due to: %s", filePath, err)
		}
	}()

	if err = processSheet(f, "BUY", BUY_TABLE_LEGEND, newBuyItem); err != nil {
		log.Error(err)
	}
	if err = processSheet(f, "SELL", SELL_TABLE_LEGEND, newSellItem); err != nil {
		log.Error(err)
	}
	if err = processSheet(f, "DIVIDEND", DIVIDEND_TABLE_LEGEND, newDividendItem); err != nil {
		log.Error(err)
	}

	return nil
}

func processSheet(excelFile *excel.File, sheetName string, legend map[string]int, newItemFunction newTransactionItem) (err error) {
	rows, err := excelFile.GetRows(sheetName, excel.Options{RawCellValue: true})
	if err != nil {
		return fmt.Errorf("sheet '%s': %s", sheetName, err)
	}
	processedRecords := 0
	for rowNo, row := range rows {
		if rowNo == 0 {
			if err := util.ValidateTableHeader(row, legend); err != nil {
				return fmt.Errorf("sheet '%s' at row '%d': %s", sheetName, rowNo, err)
			}
			continue
		} else if util.IsRowEmpty(row, legend["STOCK"]) {
			log.Debugf("Sheet '%s' at row '%d' - recognized as empty, skipping", sheetName, rowNo)
			continue
		}

		item, err := newItemFunction(row)
		if err != nil {
			return fmt.Errorf("sheet '%s' at row '%d': %s", sheetName, rowNo, err)
		}

		log.Println(item)
		processedRecords++
	}
	if processedRecords == 0 {
		log.Warnf("Sheet '%s' has not data to process", sheetName)
	}
	return nil
}
