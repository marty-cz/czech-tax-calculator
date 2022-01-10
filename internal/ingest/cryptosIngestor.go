package ingest

import (
	"fmt"
	"strconv"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
	excel "github.com/xuri/excelize/v2"
)

var (
	CRYPTO_BUY_TBL_LEGEND = map[string]int{
		"CRYPTO":     0,
		"DATE":       1,
		"COIN PRICE": 2,
		"PAID":       3,
		"FEE":        4,
		"AMOUNT":     5,
		"QUANTITY":   6,
		"BROKER":     7,
		"CURRENCY":   8,
	}
	CRYPTO_SELL_TBL_LEGEND = map[string]int{
		"CRYPTO":     0,
		"DATE":       1,
		"COIN PRICE": 2,
		"RECEIVED":   3,
		"FEE":        4,
		"AMOUNT":     5,
		"QUANTITY":   6,
		"BROKER":     7,
		"CURRENCY":   8,
	}
)

func newCryptoBuyItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[CRYPTO_BUY_TBL_LEGEND["CRYPTO"]],
		Broker:    row[CRYPTO_BUY_TBL_LEGEND["BROKER"]],
		Operation: BUY,
	}

	if rawDate, err := strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["COIN PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("coin price is not a number: %s", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["PAID"]], 64); err != nil {
		return nil, fmt.Errorf("paid is not a number: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[CRYPTO_BUY_TBL_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[CRYPTO_BUY_TBL_LEGEND["CURRENCY"]]); err != nil {
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

func newCryptoSellItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[CRYPTO_SELL_TBL_LEGEND["CRYPTO"]],
		Broker:    row[CRYPTO_SELL_TBL_LEGEND["BROKER"]],
		Operation: SELL,
	}

	if rawDate, err := strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["COIN PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("coin price is not a number: %s", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[CRYPTO_SELL_TBL_LEGEND["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[CRYPTO_SELL_TBL_LEGEND["CURRENCY"]]); err != nil {
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

func ProcessCryptos(filePath string) (_ *TransactionLog, err error) {

	log.Infof("Processing cryptos input file '%s'", filePath)

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

	log.Infof("Cryptos: Ingesting Purchases")
	if transactions.Purchases, err = processSheet(f, "BUY", CRYPTO_BUY_TBL_LEGEND, newCryptoBuyItem); err != nil {
		log.Error(err)
	}
	log.Infof("Cryptos: Ingesting Sales")
	if transactions.Sales, err = processSheet(f, "SELL", CRYPTO_SELL_TBL_LEGEND, newCryptoSellItem); err != nil {
		log.Error(err)
	}
	log.Infof("Cryptos: Ingesting Additional Income")
	if transactions.AdditionalIncomes, err = processSheet(f, "ADDITIONAL INCOME", ADDITIONAL_INCOME_TBL_LEGEND, newAdditionalIncomeItem); err != nil {
		log.Error(err)
	}
	log.Infof("Cryptos: Ingesting Additional Fee")
	if transactions.AdditionalFees, err = processSheet(f, "ADDITIONAL FEE", ADDITIONAL_FEE_TBL_LEGEND, newAdditionalFeeItem); err != nil {
		log.Error(err)
	}

	return &transactions, nil
}
