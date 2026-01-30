package ingest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
	excel "github.com/xuri/excelize/v2"
)

const MaxAllowedTax float64 = 0.15
const StockItemType string = "stock"

var (
	stockBuyTblLegend = map[string]int{
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
	stockSellTblLegend = map[string]int{
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
	stockDividendTblLegend = map[string]int{
		"STOCK":    0,
		"DATE":     1,
		"RECEIVED": 2,
		"AMOUNT":   3,
		"PAID TAX": 4,
		"BROKER":   5,
		"CURRENCY": 6,
		"COUNTRY":  7,
	}
)

func newStockBuyItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[stockBuyTblLegend["STOCK"]],
		Broker:    row[stockBuyTblLegend["BROKER"]],
		Operation: BUY,
	}

	if rawDate, err := strconv.ParseFloat(row[stockBuyTblLegend["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[stockBuyTblLegend["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %v", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[stockBuyTblLegend["PAID"]], 64); err != nil {
		return nil, fmt.Errorf("paid is not a number: %v", err)
	}
	item.OriginalBankAmount = item.BankAmount
	if item.BrokerAmount, err = strconv.ParseFloat(row[stockBuyTblLegend["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[stockBuyTblLegend["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %v", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[stockBuyTblLegend["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[stockBuyTblLegend["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %v", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get day exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	return validateStockBuyItem(&item)
}

func newStockSellItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[stockSellTblLegend["STOCK"]],
		Broker:    row[stockSellTblLegend["BROKER"]],
		Operation: SELL,
	}

	if rawDate, err := strconv.ParseFloat(row[stockSellTblLegend["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[stockSellTblLegend["STOCK PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("stock price is not a number: %v", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[stockSellTblLegend["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %v", err)
	}
	item.OriginalBankAmount = item.BankAmount
	if item.BrokerAmount, err = strconv.ParseFloat(row[stockSellTblLegend["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[stockSellTblLegend["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %v", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[stockSellTblLegend["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[stockSellTblLegend["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %v", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	return validateStockSellItem(&item)
}

func newStockDividendItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[stockDividendTblLegend["STOCK"]],
		Broker:    row[stockDividendTblLegend["BROKER"]],
		Operation: DIVIDEND,
	}

	if rawDate, err := strconv.ParseFloat(row[stockDividendTblLegend["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[stockDividendTblLegend["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %v", err)
	}
	item.OriginalBankAmount = item.BankAmount
	if item.BrokerAmount, err = strconv.ParseFloat(row[stockDividendTblLegend["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[stockDividendTblLegend["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %v", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	item.Country = strings.ToUpper(row[stockDividendTblLegend["COUNTRY"]])
	if item.Country == "" {
		return nil, fmt.Errorf("cannot get country")
	}
	return validateDividendItem(&item)
}


func validateStockBuyItem(item *TransactionLogItem) (_ *TransactionLogItem, err error) {
	if !util.LeqWithTolerance(item.BrokerAmount, item.BankAmount, 0.0001) {
		return nil, fmt.Errorf("Bank amount (PAID) is greater than Broker amount (AMOUNT) for item '%v'", item)
	}
	if !util.EqWithTolerance(item.BrokerAmount + item.Fee, item.BankAmount, 0.0001) {
		return nil, fmt.Errorf("Bank amount (PAID) is not equal to Broker amount (AMOUNT) + Fee for item '%v'", item)
	}
	return item, nil
}

func validateStockSellItem(item *TransactionLogItem) (_ *TransactionLogItem, err error) {
	if !util.LeqWithTolerance(item.BankAmount, item.BrokerAmount, 0.0001) {
		return nil, fmt.Errorf("Bank amount (RECEIVED) is greater than Broker amount (AMOUNT) for item '%v'", item)
	}
	if !util.EqWithTolerance(item.BankAmount, item.BrokerAmount - item.Fee, 0.0001) {
		return nil, fmt.Errorf("Bank amount (PAID) is not equal to Broker amount (AMOUNT) - Fee for item '%v'", item)
	}
	return item, nil
}


func validateDividendItem(item *TransactionLogItem) (_ *TransactionLogItem, err error) {
	if !util.LeqWithTolerance(item.BankAmount, item.BrokerAmount, 0.0001) {
		return nil, fmt.Errorf("Bank amount (RECEIVED) is greater than Broker amount (AMOUNT) for item '%v'", item)
	}
	paidTax := 1 - (item.BankAmount / item.BrokerAmount)
	if !util.LeqWithTolerance(paidTax, MaxAllowedTax, 0.01) {
		item.BankAmount = item.BrokerAmount * (1 - MaxAllowedTax)
		log.Warnf("Paid tax '%f' exceeds max allowed tax '%v' for item '%v' - adjusting Bank Amount to %f", paidTax, MaxAllowedTax, item, item.BankAmount)
	}	
	return item, nil
}

func ProcessStocks(filePath string) (_ *TransactionLog, err error) {
	log.Infof("%ss: processing input file '%s'", StockItemType, filePath)

	f, err := excel.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet
		if err := f.Close(); err != nil {
			log.Errorf("cannot close file '%s' due to: %v", filePath, err)
		}
	}()

	transactions := TransactionLog{}

	log.Infof("%ss: Ingesting Purchases", StockItemType)
	if transactions.Purchases, err = processSheet(f, "BUY", stockBuyTblLegend, newStockBuyItem); err != nil {
		log.Errorf("%ss: %v", StockItemType, err)
	}
	log.Infof("%ss: Ingested Purchases (count: %d)", StockItemType, len(transactions.Purchases))

	log.Infof("%ss: Ingesting Sales", StockItemType)
	if transactions.Sales, err = processSheet(f, "SELL", stockSellTblLegend, newStockSellItem); err != nil {
		log.Errorf("%ss: %v", StockItemType, err)
	}
	log.Infof("%ss: Ingested Sales (count: %d)", StockItemType, len(transactions.Sales))

	log.Infof("%ss: Ingesting Dividends", StockItemType)
	if transactions.Dividends, err = processSheet(f, "DIVIDEND", stockDividendTblLegend, newStockDividendItem); err != nil {
		log.Errorf("%ss: %v", StockItemType, err)
	}
	log.Infof("%ss: Ingested Dividends (count: %d)", StockItemType, len(transactions.Dividends))

	log.Infof("%ss: Ingesting Additional Incomes", StockItemType)
	if transactions.AdditionalIncomes, err = processSheet(f, "ADDITIONAL INCOME", ADDITIONAL_INCOME_TBL_LEGEND, newAdditionalIncomeItem); err != nil {
		log.Errorf("%ss: %v", StockItemType, err)
	}
	log.Infof("%ss: Ingested Additional Incomes (count: %d)", StockItemType, len(transactions.AdditionalIncomes))

	log.Infof("%ss: Ingesting Additional Fees", StockItemType)
	if transactions.AdditionalFees, err = processSheet(f, "ADDITIONAL FEE", ADDITIONAL_FEE_TBL_LEGEND, newAdditionalFeeItem); err != nil {
		log.Errorf("%ss: %v", StockItemType, err)
	}
	log.Infof("%ss: Ingested Additional Fees (count: %d)", StockItemType, len(transactions.AdditionalFees))

	return &transactions, nil
}
