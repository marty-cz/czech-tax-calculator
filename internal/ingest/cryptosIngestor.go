package ingest

import (
	"fmt"
	"strconv"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
	excel "github.com/xuri/excelize/v2"
)

const CryptoItemType string = "crypto"

var (
	cryptoBuyTblLegend = map[string]int{
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
	cryptoSellTblLegend = map[string]int{
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
		Name:      row[cryptoBuyTblLegend["CRYPTO"]],
		Broker:    row[cryptoBuyTblLegend["BROKER"]],
		Operation: BUY,
	}

	if rawDate, err := strconv.ParseFloat(row[cryptoBuyTblLegend["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[cryptoBuyTblLegend["COIN PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("coin price is not a number: %v", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[cryptoBuyTblLegend["PAID"]], 64); err != nil {
		return nil, fmt.Errorf("paid is not a number: %v", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[cryptoBuyTblLegend["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[cryptoBuyTblLegend["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %v", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[cryptoBuyTblLegend["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[cryptoBuyTblLegend["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %v", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get day exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	return &item, nil
}

func newCryptoSellItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      row[cryptoSellTblLegend["CRYPTO"]],
		Broker:    row[cryptoSellTblLegend["BROKER"]],
		Operation: SELL,
	}

	if rawDate, err := strconv.ParseFloat(row[cryptoSellTblLegend["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.ItemPrice, err = strconv.ParseFloat(row[cryptoSellTblLegend["COIN PRICE"]], 64); err != nil {
		return nil, fmt.Errorf("coin price is not a number: %v", err)
	}
	if item.BankAmount, err = strconv.ParseFloat(row[cryptoSellTblLegend["RECEIVED"]], 64); err != nil {
		return nil, fmt.Errorf("received is not a number: %v", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[cryptoSellTblLegend["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[cryptoSellTblLegend["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %v", err)
	}
	if item.Quantity, err = strconv.ParseFloat(row[cryptoSellTblLegend["QUANTITY"]], 64); err != nil {
		return nil, fmt.Errorf("quantity is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[cryptoSellTblLegend["CURRENCY"]]); err != nil {
		return nil, fmt.Errorf("currency format problem: %v", err)
	}
	if item.DayExchangeRate, err = util.GetCzkExchangeRateInDay(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	if item.YearExchangeRate, err = util.GetCzkExchangeRateInYear(item.Date, *item.Currency); err != nil {
		return nil, fmt.Errorf("cannot get year exchange rate for %v from %v: %v", item.Currency, item.Date, err)
	}
	return &item, nil
}

func ProcessCryptos(filePath string) (_ *TransactionLog, err error) {

	log.Infof("%ss: processing input file '%s'", CryptoItemType, filePath)

	f, err := excel.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Errorf("cannot close file '%s' due to: %v", filePath, err)
		}
	}()

	transactions := TransactionLog{}

	log.Infof("%ss: Ingesting Purchases", CryptoItemType)
	if transactions.Purchases, err = processSheet(f, "BUY", cryptoBuyTblLegend, newCryptoBuyItem); err != nil {
		log.Errorf("%ss: %v", CryptoItemType, err)
	}
	log.Infof("%ss: Ingested Purchases (count: %d)", CryptoItemType, len(transactions.Purchases))

	log.Infof("%ss: Ingesting Sales", CryptoItemType)
	if transactions.Sales, err = processSheet(f, "SELL", cryptoSellTblLegend, newCryptoSellItem); err != nil {
		log.Errorf("%ss: %v", CryptoItemType, err)
	}
	log.Infof("%ss: Ingested Sales (count: %d)", CryptoItemType, len(transactions.Sales))

	log.Infof("%ss: Ingesting Additional Incomes", CryptoItemType)
	if transactions.AdditionalIncomes, err = processSheet(f, "ADDITIONAL INCOME", ADDITIONAL_INCOME_TBL_LEGEND, newAdditionalIncomeItem); err != nil {
		log.Errorf("%ss: %v", CryptoItemType, err)
	}
	log.Infof("%ss: Ingested Additional Incomes (count: %d)", CryptoItemType, len(transactions.AdditionalIncomes))

	log.Infof("%ss: Ingesting Additional Fees", CryptoItemType)
	if transactions.AdditionalFees, err = processSheet(f, "ADDITIONAL FEE", ADDITIONAL_FEE_TBL_LEGEND, newAdditionalFeeItem); err != nil {
		log.Errorf("%ss: %v", CryptoItemType, err)
	}
	log.Infof("%ss: Ingested Additional Fees (count: %d)", CryptoItemType, len(transactions.AdditionalFees))

	return &transactions, nil
}
