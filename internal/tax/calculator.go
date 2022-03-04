package tax

import (
	"fmt"
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

	itemsToSell := convertToItemsToSell(transactions.Purchases)

	// go through tax years from oldest to latest
	for year := oldestSellTransactionYear; year <= currentTaxYear; year++ {
		inYearSellOperations, dateStart, dateEnd, err := getItemSales(transactions.Sales, itemsToSell, year, allowThreeYearsTimeTest)
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

func getItemSales(sellTransactions ingest.TransactionLogItems, itemsToSell ItemsToSell, year int, allowThreeYearsTimeTest bool) (SellOperations, time.Time, time.Time, error) {
	layout := "02.01.2006 15:04:05"
	dateStart, _ := time.Parse(layout, fmt.Sprintf("01.01.%d 00:00:00", year))
	dateEnd, _ := time.Parse(layout, fmt.Sprintf("31.12.%d 23:59:59", year))

	inYearSellTransactions := getTransactionsInYear(sellTransactions, dateStart, dateEnd)
	inYearSellOperations := convertToSellOperations(inYearSellTransactions)

	log.Infof("sale transactions count for year '%d': %d", year, len(inYearSellOperations))

	for _, sellOp := range inYearSellOperations {
		availableBuyItems := getAvailableItemsToSell(itemsToSell, sellOp.SellItem)
		log.Debugf("sell '%s' available buy items: %v", sellOp.SellItem.Name, availableBuyItems)

		calculateSellExpense(sellOp, availableBuyItems, allowThreeYearsTimeTest)
		log.Debugf("sell operation processed: '%+v'", sellOp)
	}
	return inYearSellOperations, dateStart, dateEnd, nil
}

func calculateReport(sellOps SellOperations, dividends ingest.TransactionLogItems, additionalIncomes ingest.TransactionLogItems, additionalFees ingest.TransactionLogItems, year time.Time) *Report {
	report := Report{
		SellOperations:                sellOps,
		Year:                          year,
		Currency:                      DEFAULT_CURRENCY,
		TotalItemRevenue:              newAccountingValue(0, 0, DEFAULT_CURRENCY),
		TimeTestedItemRevenue:         newAccountingValue(0, 0, DEFAULT_CURRENCY),
		DividendRevenue:               newEmptyValueAndFee(DEFAULT_CURRENCY),
		DividendRevenueWithNoTaxPayed: newEmptyValueAndFee(DEFAULT_CURRENCY),
		AdditionalRevenue:             newEmptyValueAndFee(DEFAULT_CURRENCY),
		TimeTestedItemFifoExpense:     newEmptyValueAndFee(DEFAULT_CURRENCY),
		TotalItemFifoExpense:          newEmptyValueAndFee(DEFAULT_CURRENCY),
	}

	// calculate report for sold items
	for _, sellOp := range sellOps {
		timeTestedQuantity := 0.0
		// sum expenses and buy fees of all sold items
		for _, soldItem := range sellOp.SoldItems {
			report.TotalItemFifoExpense.Value.Add(soldItem.FifoBuy.Value)
			report.TotalItemFifoExpense.Fee.Add(soldItem.FifoBuy.Fee)
			if soldItem.TimeTested {
				timeTestedQuantity += soldItem.SoldQuantity
				report.TimeTestedItemFifoExpense.Value.Add(soldItem.FifoBuy.Value)
				report.TimeTestedItemFifoExpense.Fee.Add(soldItem.FifoBuy.Fee)
			}
		}
		timeTestedRatio := timeTestedQuantity / sellOp.SellItem.Quantity

		// calculate revenue and sell fees
		report.TotalItemRevenue.Add(sellOp.totalRevenue)
		sellStockFee := newAccountingValue(
			sellOp.SellItem.Fee*sellOp.SellItem.DayExchangeRate,
			sellOp.SellItem.Fee*sellOp.SellItem.YearExchangeRate,
			report.Currency)
		report.TotalItemFifoExpense.Fee.Add(sellStockFee)
		report.TimeTestedItemRevenue.Add(sellOp.timeTestedRevenue)
		report.TimeTestedItemFifoExpense.Fee.Add(sellStockFee.MultiplyNew(timeTestedRatio))
	}

	// calculate report for received dividends
	for _, dividend := range dividends {
		divRevenue := *report.DividendRevenueWithNoTaxPayed
		if dividend.BankAmount < dividend.BrokerAmount {
			divRevenue = *report.DividendRevenue
		}
		divRevenue.Value.Add(newAccountingValue(
			dividend.BrokerAmount*dividend.DayExchangeRate,
			dividend.BrokerAmount*dividend.YearExchangeRate, report.Currency))
		divRevenue.Fee.Add(newAccountingValue(
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

func calculateSellExpense(sellOp *SellOperation, availableBuyItems ItemsToSell, allowThreeYearsTimeTest bool) {

	timeTestDate := util.GetDateThreeYearsBefore(sellOp.SellItem.Date)

	quantityToBeSold := sellOp.SellItem.Quantity
	for _, itemToSell := range availableBuyItems {
		if itemToSell.availableQuantity <= 0.0 {
			continue
		}

		soldItem := &SoldItem{
			BuyItem: itemToSell.buyItem,
			// 3 years time test
			TimeTested: allowThreeYearsTimeTest && itemToSell.buyItem.Date.Sub(timeTestDate).Nanoseconds() < 0,
			FifoBuy:    newEmptyValueAndFee(DEFAULT_CURRENCY),
			Revenue:    newEmptyValueAndFee(DEFAULT_CURRENCY),
		}
		soldBuyItemRatio := 0.0
		newAvailableQuantity := itemToSell.availableQuantity - quantityToBeSold
		if newAvailableQuantity >= 0.0 {
			// sell operation has all buys processed
			soldItem.SoldQuantity = quantityToBeSold
			itemToSell.availableQuantity = newAvailableQuantity
			soldBuyItemRatio = soldItem.SoldQuantity / itemToSell.buyItem.Quantity
		} else {
			// some buy items are still required to be sold by this sell operation
			quantityToBeSold -= itemToSell.availableQuantity
			soldItem.SoldQuantity = itemToSell.availableQuantity
			// noting remains to be sold in the buy item
			itemToSell.availableQuantity = 0.0
			soldBuyItemRatio = soldItem.SoldQuantity / itemToSell.buyItem.Quantity
		}

		// calculate purchase for this item
		soldItem.FifoBuy.Value = newAccountingValue(
			soldBuyItemRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.DayExchangeRate,
			soldBuyItemRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		soldItem.FifoBuy.Fee = newAccountingValue(
			soldBuyItemRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.DayExchangeRate,
			soldBuyItemRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		// calculate revenue for this buy item
		soldRatio := soldItem.SoldQuantity / sellOp.SellItem.Quantity
		soldItem.Revenue.Value = newAccountingValue(
			soldRatio*sellOp.SellItem.BrokerAmount*sellOp.SellItem.DayExchangeRate,
			soldRatio*sellOp.SellItem.BrokerAmount*sellOp.SellItem.YearExchangeRate,
			DEFAULT_CURRENCY)
		soldItem.Revenue.Fee = newAccountingValue(
			soldRatio*sellOp.SellItem.Fee*sellOp.SellItem.DayExchangeRate,
			soldRatio*sellOp.SellItem.Fee*sellOp.SellItem.YearExchangeRate,
			DEFAULT_CURRENCY)

		if soldItem.TimeTested {
			sellOp.timeTestedRevenue.Add(soldItem.Revenue.Value)
		}
		sellOp.totalRevenue.Add(soldItem.Revenue.Value)

		sellOp.SoldItems = append(sellOp.SoldItems, soldItem)
		itemToSell.soldByItems = append(itemToSell.soldByItems, sellOp.SellItem)

		if newAvailableQuantity >= 0.0 {
			return
		}
	}
}

func getTransactionsInYear(transactions ingest.TransactionLogItems, from time.Time, to time.Time) (ret ingest.TransactionLogItems) {
	fromExclusive := from.Add(-1 * time.Second)
	toExclusive := to.Add(1 * time.Second)
	isTransactionTimestampBetween := func(item *ingest.TransactionLogItem) bool {
		return item.Date.After(fromExclusive) && item.Date.Before(toExclusive)
	}
	return filterTransactionLog(transactions, isTransactionTimestampBetween)
}

func filterTransactionLog(list ingest.TransactionLogItems, test func(*ingest.TransactionLogItem) bool) (ret ingest.TransactionLogItems) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}
