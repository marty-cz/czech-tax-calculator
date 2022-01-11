package ingest

import (
	"fmt"
	"strconv"

	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
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
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.BrokerAmount, err = strconv.ParseFloat(row[ADDITIONAL_INCOME_TBL_LEGEND["AMOUNT"]], 64); err != nil {
		return nil, fmt.Errorf("amount is not a number: %v", err)
	}
	item.BankAmount = item.BrokerAmount
	if item.Currency, err = util.GetCurrencyByName(row[ADDITIONAL_INCOME_TBL_LEGEND["CURRENCY"]]); err != nil {
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

func newAdditionalFeeItem(row []string) (_ *TransactionLogItem, err error) {
	item := TransactionLogItem{
		Name:      "",
		Broker:    row[ADDITIONAL_FEE_TBL_LEGEND["LOCATION"]],
		Operation: ADDITIONAL_FEE,
	}

	if rawDate, err := strconv.ParseFloat(row[ADDITIONAL_FEE_TBL_LEGEND["DATE"]], 64); err != nil {
		return nil, fmt.Errorf("raw date is not a number: %v", err)
	} else if item.Date, err = excel.ExcelDateToTime(rawDate, false); err != nil {
		return nil, fmt.Errorf("date has invalid format: %v", err)
	}
	if item.Fee, err = strconv.ParseFloat(row[ADDITIONAL_FEE_TBL_LEGEND["FEE"]], 64); err != nil {
		return nil, fmt.Errorf("fee is not a number: %v", err)
	}
	if item.Currency, err = util.GetCurrencyByName(row[ADDITIONAL_FEE_TBL_LEGEND["CURRENCY"]]); err != nil {
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

func processSheet(excelFile *excel.File, sheetName string, legend map[string]int, newItemFunction newTransactionItem) (transactions TransactionLogItems, err error) {
	rows, err := excelFile.GetRows(sheetName, excel.Options{RawCellValue: true})
	if err != nil {
		return nil, fmt.Errorf("sheet '%s': %v", sheetName, err)
	}

	transactions = make(TransactionLogItems, 0)
	for rowNo, row := range rows {
		excelRowNo := rowNo + 1
		if rowNo == 0 {
			if err := util.ValidateTableHeader(row, legend); err != nil {
				return nil, fmt.Errorf("sheet '%s' (row '%d'): %v", sheetName, rowNo, err)
			}
			continue
		} else if util.IsRowEmpty(row, 0) {
			log.Warnf("sheet '%s' (row '%d') - recognized as empty, skipping", sheetName, excelRowNo)
			continue
		}

		item, err := newItemFunction(row)
		if err != nil {
			return nil, fmt.Errorf("sheet '%s' (row '%d'): %v", sheetName, excelRowNo, err)
		}

		transactions = append(transactions, item)
		log.Debugf("ingested from '%s' (row '%d'): %+v", sheetName, excelRowNo, item)
	}
	if len(transactions) == 0 {
		log.Warnf("sheet '%s' has not data to process", sheetName)
	}
	return transactions, nil
}
