package ingest

import (
	"fmt"
	"strconv"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	excel "github.com/xuri/excelize/v2"
)

var (
	ADDITIONAL_INCOME_TBL_LEGEND = map[string]int{
		"DATE":     0,
		"AMOUNT":   1,
		"LOCATION": 2,
		"CURRENCY": 3,
	}
	ADDITIONAL_FEE_TBL_LEGEND = map[string]int{
		"DATE":     0,
		"FEE":      1,
		"LOCATION": 2,
		"CURRENCY": 3,
	}
)

func newAdditionalIncomeItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      "",
		Broker:    row[ADDITIONAL_INCOME_TBL_LEGEND["LOCATION"]],
		Operation: ADDITIONAL_INCOME,
	}

	if rawDate, err := strconv.ParseFloat(row[ADDITIONAL_INCOME_TBL_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[ADDITIONAL_INCOME_TBL_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %s", err)
	}
	item.BankAmount = item.BrokerAmount
	if item.Currency, err = util.GetCurrencyByName(row[ADDITIONAL_INCOME_TBL_LEGEND["CURRENCY"]]); err != nil {
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

func newAdditionalFeeItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      "",
		Broker:    row[ADDITIONAL_FEE_TBL_LEGEND["LOCATION"]],
		Operation: ADDITIONAL_FEE,
	}

	if rawDate, err := strconv.ParseFloat(row[ADDITIONAL_FEE_TBL_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %s", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %s", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[ADDITIONAL_FEE_TBL_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %s", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[ADDITIONAL_FEE_TBL_LEGEND["CURRENCY"]]); err != nil {
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
