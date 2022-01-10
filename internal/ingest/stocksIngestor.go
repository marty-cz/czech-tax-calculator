package ingest

import (
	"fmt"
	"strconv"

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

func newBuyItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[BUY_TABLE_LEGEND["STOCK"]],
		Broker:    row[BUY_TABLE_LEGEND["BROKER"]],
		Operation: BUY,
	}

	if rawDate, err := strconv.ParseFloat(row[BUY_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %s", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["PAID"]], 64); err != nil {
		return nil, fmt.Errorf("paid is not a number: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[BUY_TABLE_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[BUY_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get day exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	return &item, nil
}

func newSellItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[SELL_TABLE_LEGEND["STOCK"]],
		Broker:    row[SELL_TABLE_LEGEND["BROKER"]],
		Operation: SELL,
	}

	if rawDate, err := strconv.ParseFloat(row[SELL_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %s", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[SELL_TABLE_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[SELL_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	return &item, nil
}

func newDividendItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[DIVIDEND_TABLE_LEGEND["STOCK"]],
		Broker:    row[DIVIDEND_TABLE_LEGEND["BROKER"]],
		Operation: DIVIDEND,
	}

	if rawDate, err := strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[DIVIDEND_TABLE_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[DIVIDEND_TABLE_LEGEND["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %s", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %s", item.Currency, item.Date, err)
	}
	return &item, nil
}

func ProcessStocks(filePath string) (_ *TransactionLog, err error) {

	log.Infof("Processing stocks input file '%s'", filePath)

	f, err := excel.OpenFile(filePath)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Errorf("Cannot close file '%s' due to: %s", filePath, err)
		}
	}()

	transactions := TransactionLog{}

	log.Infof("Stocks: Ingesting Purchases")
	if transactions.Purchases, err = processSheet(f, "BUY", BUY_TABLE_LEGEND, newBuyItem); err != nil {
		log.Error(err)
	}
	log.Infof("Stocks: Ingesting Sales")
	if transactions.Sales, err = processSheet(f, "SELL", SELL_TABLE_LEGEND, newSellItem); err != nil {
		log.Error(err)
	}
	log.Infof("Stocks: Ingesting Dividends")
	if transactions.Dividends, err = processSheet(f, "DIVIDEND", DIVIDEND_TABLE_LEGEND, newDividendItem); err != nil {
		log.Error(err)
	}
	log.Infof("Stocks: Ingesting Additional Income")
	if transactions.AdditionalIncomes, err = processSheet(f, "ADDITIONAL INCOME", ADDITIONAL_INCOME_TABLE_LEGEND, newAdditionalIncomeItem); err != nil {
		log.Error(err)
	}
	log.Infof("Stocks: Ingesting Additional Fee")
	if transactions.AdditionalFees, err = processSheet(f, "ADDITIONAL FEE", ADDITIONAL_FEE_TABLE_LEGEND, newAdditionalFeeItem); err != nil {
		log.Error(err)
	}

	return &transactions, nil
}
