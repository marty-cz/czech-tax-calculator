package tax

import (
	"fmt"
	"strings"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
)

var DEFAULT_CURRENCY *util.Currency = util.CZK

func Calculate(transactions *ingest.TransactionLog, currentTaxYearString string, allowThreeYearsTimeTest bool) (reports Reports, err error) {
	currentTaxYear, err := util.GetYearFromString(currentTaxYearString)
	if err != nil {
		return nil, err
	}

	// find year of oldest revenue transaction
	sortByDate(transactions)
	oldestSellTransactionYear := currentTaxYear
	if len(transactions.Sales) > 0 {
		oldestSellTransactionYear = transactions.Sales[0].Date.Year()
	}
	if len(transactions.AdditionalIncomes) > 0 && transactions.AdditionalIncomes[0].Date.Year() < oldestSellTransactionYear {
		oldestSellTransactionYear = transactions.AdditionalIncomes[0].Date.Year()
	}

	// go through tax years from oldest to latest
	for year := oldestSellTransactionYear; year <= currentTaxYear; year++ {
		inYearSellOperations, dateStart, dateEnd, err := getItemSales(transactions, year, allowThreeYearsTimeTest)
		if err != nil {
			return nil, fmt.Errorf("calculation for year '%v' failed: %v", year, err)
		}
		inYearDividends := getTransactionsInYear(transactions.Dividends, dateStart, dateEnd)
		inYearAdditionalIncomes := getTransactionsInYear(transactions.AdditionalIncomes, dateStart, dateEnd)
		inYearAdditionalFees := getTransactionsInYear(transactions.AdditionalFees, dateStart, dateEnd)
		reports = append(reports, calculateReport(inYearSellOperations, inYearDividends, inYearAdditionalIncomes, inYearAdditionalFees, dateStart))
	}

	return
}

func getItemSales(transactions *ingest.TransactionLog, year int, allowThreeYearsTimeTest bool) (SellOperations, time.Time, time.Time, error) {
	layout := "02.01.2006 15:04:05"
	dateStart, _ := time.Parse(layout, fmt.Sprintf("01.01.%d 00:00:00", year))
	dateEnd, _ := time.Parse(layout, fmt.Sprintf("31.12.%d 23:59:59", year))

	itemsToSell := convertToItemsToSell(transactions.Purchases)
	sellOperations := convertToSellOperations(transactions.Sales)

	inYearSellOperations := getSalesInYear(sellOperations, dateStart, dateEnd)
	log.Infof("sale transactions count for year '%d': %d", year, len(inYearSellOperations))

	for _, sellOp := range inYearSellOperations {
		availableBuyItems := getAvailableItemsToSell(itemsToSell, sellOp.sellItem)
		log.Debugf("sell '%s' available buy items: %v", sellOp.sellItem.Name, availableBuyItems)

		calculateSellExpense(sellOp, availableBuyItems, allowThreeYearsTimeTest)
		log.Debugf("sell operation processed: '%+v'", sellOp)
	}
	return inYearSellOperations, dateStart, dateEnd, nil
}

func calculateReport(sellOps SellOperations, dividends ingest.TransactionLogItems, additionalIncomes ingest.TransactionLogItems, additionalFees ingest.TransactionLogItems, year time.Time) *Report {
	report := Report{
		SellOperations:            sellOps,
		Year:                      year,
		Currency:                  DEFAULT_CURRENCY,
		TotalItemRevenue:          newAccountingValue(0, 0, DEFAULT_CURRENCY),
		TimeTestedItemRevenue:     newAccountingValue(0, 0, DEFAULT_CURRENCY),
		DividendRevenue:           newEmptyValueAndFee(DEFAULT_CURRENCY),
		AdditionalRevenue:         newEmptyValueAndFee(DEFAULT_CURRENCY),
		TimeTestedItemFifoExpense: newEmptyValueAndFee(DEFAULT_CURRENCY),
		TotalItemFifoExpense:      newEmptyValueAndFee(DEFAULT_CURRENCY),
	}

	// calculate report for sold items
	for _, sellOp := range sellOps {
		timeTestedQuantity := 0.0
		// sum expenses and buy fees of all sold items
		for _, soldItem := range sellOp.soldItems {
			report.TotalItemFifoExpense.Value.Add(soldItem.fifoBuy.Value)
			report.TotalItemFifoExpense.Fee.Add(soldItem.fifoBuy.Fee)
			if soldItem.timeTested {
				timeTestedQuantity += soldItem.soldQuantity
				report.TimeTestedItemFifoExpense.Value.Add(soldItem.fifoBuy.Value)
				report.TimeTestedItemFifoExpense.Fee.Add(soldItem.fifoBuy.Fee)
			}
		}
		timeTestedRatio := timeTestedQuantity / sellOp.sellItem.Quantity

		// calculate revenue and sell fees
		report.TotalItemRevenue.Add(sellOp.totalRevenue)
		sellStockFee := newAccountingValue(
			sellOp.sellItem.Fee*sellOp.sellItem.DayExchangeRate,
			sellOp.sellItem.Fee*sellOp.sellItem.YearExchangeRate,
			report.Currency)
		report.TotalItemFifoExpense.Fee.Add(sellStockFee)
		report.TimeTestedItemRevenue.Add(sellOp.timeTestedRevenue)
		report.TimeTestedItemFifoExpense.Fee.Add(sellStockFee.MultiplyNew(timeTestedRatio))
	}

	// calculate report for received dividends
	for _, dividend := range dividends {
		report.DividendRevenue.Value.Add(newAccountingValue(
			dividend.BrokerAmount*dividend.DayExchangeRate,
			dividend.BrokerAmount*dividend.YearExchangeRate, report.Currency))
		report.DividendRevenue.Fee.Add(newAccountingValue(
			dividend.Fee*dividend.DayExchangeRate,
			dividend.Fee*dividend.YearExchangeRate, report.Currency))
	}
	// calculate report for received additional income
	for _, additionalIncome := range additionalIncomes {
		report.AdditionalRevenue.Value.Add(newAccountingValue(
			additionalIncome.BrokerAmount*additionalIncome.DayExchangeRate,
			additionalIncome.BrokerAmount*additionalIncome.YearExchangeRate, report.Currency))
	}
	// calculate report for paid additional fees
	for _, additionalFee := range additionalFees {
		report.AdditionalRevenue.Fee.Add(newAccountingValue(
			additionalFee.BrokerAmount*additionalFee.DayExchangeRate,
			additionalFee.BrokerAmount*additionalFee.YearExchangeRate, report.Currency))
	}

	return &report
}

func calculateSellExpense(sellOp *SellOperation, availableBuyItems ItemToSellCollection, allowThreeYearsTimeTest bool) {

	timeTestDate := util.GetDateThreeYearsBefore(sellOp.sellItem.Date)

	quantityToBeSold := sellOp.sellItem.Quantity
	for _, itemToSell := range availableBuyItems {
		if itemToSell.availableQuantity <= 0.0 {
			continue
		}

		soldItem := &SoldItem{
			buyItem: itemToSell.buyItem,
			// 3 years time test
			timeTested: allowThreeYearsTimeTest && itemToSell.buyItem.Date.Sub(timeTestDate).Nanoseconds() < 0,
			fifoBuy:    newEmptyValueAndFee(DEFAULT_CURRENCY),
		}
		soldBuyItemRatio := 0.0
		newAvailableQuantity := itemToSell.availableQuantity - quantityToBeSold
		if newAvailableQuantity >= 0.0 {
			// sell operation has all buys processed
			soldItem.soldQuantity = quantityToBeSold
			itemToSell.availableQuantity = newAvailableQuantity
			soldBuyItemRatio = soldItem.soldQuantity / itemToSell.buyItem.Quantity
		} else {
			// some buy items are still required to be sold by this sell operation
			quantityToBeSold -= itemToSell.availableQuantity
			soldItem.soldQuantity = itemToSell.availableQuantity
			// noting remains to be sold in the buy item
			itemToSell.availableQuantity = 0.0
			soldBuyItemRatio = soldItem.soldQuantity / itemToSell.buyItem.Quantity
		}

		// calculate purchase for this item
		soldItem.fifoBuy.Value = newAccountingValue(
			soldBuyItemRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.DayExchangeRate,
			soldBuyItemRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		soldItem.fifoBuy.Fee = newAccountingValue(
			soldBuyItemRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.DayExchangeRate,
			soldBuyItemRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		// calculate revenue for this buy item
		soldRatio := soldItem.soldQuantity / sellOp.sellItem.Quantity
		revenue := newAccountingValue(
			soldRatio*sellOp.sellItem.BrokerAmount*sellOp.sellItem.DayExchangeRate,
			soldRatio*sellOp.sellItem.BrokerAmount*sellOp.sellItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		if soldItem.timeTested {
			sellOp.timeTestedRevenue.Add(revenue)
		}
		sellOp.totalRevenue.Add(revenue)

		sellOp.soldItems = append(sellOp.soldItems, soldItem)
		itemToSell.soldByItems = append(itemToSell.soldByItems, sellOp.sellItem)

		if newAvailableQuantity >= 0.0 {
			return
		}
	}
}

func getSalesInYear(sellOperations SellOperations, from time.Time, to time.Time) (ret SellOperations) {
	fromExclusive := from.Add(-1 * time.Second)
	toExclusive := to.Add(1 * time.Second)
	isTransactionTimestampBetween := func(item *SellOperation) bool {
		return item.sellItem.Date.After(fromExclusive) && item.sellItem.Date.Before(toExclusive)
	}
	return filterSellOperations(sellOperations, isTransactionTimestampBetween)
}

func getAvailableItemsToSell(itemsToSell ItemToSellCollection, sellTransaction *ingest.TransactionLogItem) (ret ItemToSellCollection) {
	test := func(itemToSell *ItemToSell) bool {
		return strings.EqualFold(itemToSell.buyItem.Name, sellTransaction.Name) && itemToSell.availableQuantity > 0.0
	}
	return filterItemsToSell(itemsToSell, test)
}

func getTransactionsInYear(transactions ingest.TransactionLogItems, from time.Time, to time.Time) (ret ingest.TransactionLogItems) {
	fromExclusive := from.Add(-1 * time.Second)
	toExclusive := to.Add(1 * time.Second)
	isTransactionTimestampBetween := func(item *ingest.TransactionLogItem) bool {
		return item.Date.After(fromExclusive) && item.Date.Before(toExclusive)
	}
	return filterTransactionLog(transactions, isTransactionTimestampBetween)
}

func filterSellOperations(list SellOperations, test func(*SellOperation) bool) (ret SellOperations) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

func filterItemsToSell(list ItemToSellCollection, test func(*ItemToSell) bool) (ret ItemToSellCollection) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

func filterTransactionLog(list ingest.TransactionLogItems, test func(*ingest.TransactionLogItem) bool) (ret ingest.TransactionLogItems) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}
