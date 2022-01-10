package tax

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
	log "github.com/sirupsen/logrus"
)

// ByAge implements sort.Interface based on the Age field.
type ByDate ingest.TransactionLogItems

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type AccountingValue struct {
	ValueWithDayExchangeRate  float64
	ValueWithYearExchangeRate float64
	Currency                  *util.Currency
}

func newAccountingValue(valueWithDayRate, valueWithYearRate float64, c *util.Currency) *AccountingValue {
	if c == nil {
		c = util.CZK
	}
	return &AccountingValue{ValueWithDayExchangeRate: valueWithDayRate, ValueWithYearExchangeRate: valueWithYearRate, Currency: c}
}

func (x *AccountingValue) String() string {
	return fmt.Sprintf("withDayExchange:%s %v withYearExchange:%s %v", x.Currency.Symbol, x.ValueWithDayExchangeRate, x.Currency.Symbol, x.ValueWithYearExchangeRate)
}

func (x *AccountingValue) Add(add *AccountingValue) {
	x.ValueWithDayExchangeRate += add.ValueWithDayExchangeRate
	x.ValueWithYearExchangeRate += add.ValueWithYearExchangeRate
}

func (x *AccountingValue) Sub(add *AccountingValue) {
	x.ValueWithDayExchangeRate -= add.ValueWithDayExchangeRate
	x.ValueWithYearExchangeRate -= add.ValueWithYearExchangeRate
}
func (x *AccountingValue) MultiplyNew(multiplicator float64) *AccountingValue {
	return &AccountingValue{
		ValueWithDayExchangeRate:  x.ValueWithDayExchangeRate * multiplicator,
		ValueWithYearExchangeRate: x.ValueWithYearExchangeRate * multiplicator,
		Currency:                  x.Currency,
	}
}

type ItemToSell struct {
	buyItem           *ingest.TransactionLogItem
	availableQuantity float64
	soldByItems       ingest.TransactionLogItems
}

func (item *ItemToSell) String() string {
	return fmt.Sprintf("buyItem:%+v availableQuantity:%v soldByItems:%+v", item.buyItem, item.availableQuantity, &item.soldByItems)
}

type ItemToSellCollection []*ItemToSell

type SoldItem struct {
	buyItem      *ingest.TransactionLogItem
	soldQuantity float64
	timeTested   bool
	fifoBuyPrice *AccountingValue
	fifoBuyFee   *AccountingValue
}

func (x *SoldItem) String() string {
	return fmt.Sprintf("buyItem:%+v soldQuantity:%v timeTested:%v fifoBuy:(%v) fifoFee:(%v)", x.buyItem, x.soldQuantity, x.timeTested, x.fifoBuyPrice, x.fifoBuyFee)
}

type SoldItemCollection []*SoldItem

type SellOperation struct {
	sellItem          *ingest.TransactionLogItem
	soldItems         SoldItemCollection
	timeTestedRevenue *AccountingValue
	totalRevenue      *AccountingValue
}

func (x *SellOperation) String() string {
	return fmt.Sprintf("sellItem:%+v totalRevenue:(%v) timeTestedRevenue:(%v) soldItems:%+v",
		x.sellItem, x.totalRevenue, x.timeTestedRevenue, &x.soldItems)
}

type SellOperationCollection []*SellOperation

type Report struct {
	SellOperations         SellOperationCollection
	TimeTestedStockRevenue *AccountingValue
	TotalStockRevenue      *AccountingValue
	DividendRevenue        *AccountingValue
	DividendFee            *AccountingValue
	TimeTestedStockExpense *AccountingValue
	TotalStockExpense      *AccountingValue
	TimeTestedStockFee     *AccountingValue
	TotalStockFee          *AccountingValue
	Year                   time.Time
	Currency               *util.Currency
}

func (x *Report) String() string {
	return fmt.Sprintf("year: %d sellOpsCount:%d stock:(total:(revenue:(%v) expense:(%v) fee:(%v)) timeTested:(revenue:(%v) expense:(%v) fee:(%v))) dividend:(revenue:(%v) fee:(%v))",
		x.Year.Year(), len(x.SellOperations), x.TotalStockRevenue, x.TotalStockExpense, x.TotalStockFee, x.TimeTestedStockRevenue, x.TimeTestedStockExpense, x.TimeTestedStockFee, x.DividendRevenue, x.DividendFee)
}

func Calculate(transactions *ingest.TransactionLog, currentTaxYearString string) (reports []*Report, err error) {

	currentTaxYear, err := util.GetYearFromString(currentTaxYearString)
	if err != nil {
		return nil, err
	}

	sortByDate(transactions)
	oldestSellTransactionYear := currentTaxYear
	if len(transactions.Sales) > 0 {
		oldestSellTransactionYear = transactions.Sales[0].Date.Year()
	}
	for year := oldestSellTransactionYear; year <= currentTaxYear; year++ {
		inYearSellOperations, dateStart, dateEnd, err := process(transactions, year)
		if err != nil {
			return nil, fmt.Errorf("calculation for year '%v' failed: %v", year, err)
		}
		inYearDividends := getDividendsInYear(transactions.Dividends, dateStart, dateEnd)
		reports = append(reports, calculateReport(inYearSellOperations, inYearDividends, dateStart))
	}

	return
}

func process(transactions *ingest.TransactionLog, year int) (SellOperationCollection, time.Time, time.Time, error) {
	layout := "02.01.2006 15:04:05"
	dateStart, _ := time.Parse(layout, fmt.Sprintf("01.01.%d 00:00:00", year))
	dateEnd, _ := time.Parse(layout, fmt.Sprintf("31.12.%d 23:59:59", year))

	itemsToSell := convertToItemsToSell(transactions.Purchases)
	sellOperations := convertToSellOperations(transactions.Sales)

	inYearSellOperations := getSalesInYear(sellOperations, dateStart, dateEnd)
	log.Debugf("Sale transaction count for year '%s' is '%d'", year, len(inYearSellOperations))
	for _, sellOp := range inYearSellOperations {
		availableBuyItems := getAvailableItemsToSell(itemsToSell, sellOp.sellItem)
		log.Debugf("For '%s' : %v", sellOp.sellItem.Name, availableBuyItems)
		calculateSellExpense(sellOp, availableBuyItems)
		log.Debugf("Processed '%+v'", sellOp)
	}
	return inYearSellOperations, dateStart, dateEnd, nil
}

func calculateReport(sellOps SellOperationCollection, dividends ingest.TransactionLogItems, year time.Time) *Report {
	report := Report{
		SellOperations:         sellOps,
		Year:                   year,
		Currency:               util.CZK,
		TotalStockRevenue:      newAccountingValue(0, 0, util.CZK),
		TimeTestedStockRevenue: newAccountingValue(0, 0, util.CZK),
		DividendRevenue:        newAccountingValue(0, 0, util.CZK),
		DividendFee:            newAccountingValue(0, 0, util.CZK),
		TotalStockExpense:      newAccountingValue(0, 0, util.CZK),
		TimeTestedStockExpense: newAccountingValue(0, 0, util.CZK),
		TotalStockFee:          newAccountingValue(0, 0, util.CZK),
		TimeTestedStockFee:     newAccountingValue(0, 0, util.CZK),
	}

	for _, sellOp := range sellOps {
		timeTestedQuantity := 0.0
		for _, soldItem := range sellOp.soldItems {
			report.TotalStockExpense.Add(soldItem.fifoBuyPrice)
			report.TotalStockFee.Add(soldItem.fifoBuyFee)
			if soldItem.timeTested {
				timeTestedQuantity += soldItem.soldQuantity
				report.TimeTestedStockExpense.Add(soldItem.fifoBuyPrice)
				report.TimeTestedStockFee.Add(soldItem.fifoBuyFee)
			}
		}
		timeTestedRatio := timeTestedQuantity / sellOp.sellItem.Quantity

		report.TotalStockRevenue.Add(sellOp.totalRevenue)
		sellStockFee := newAccountingValue(
			sellOp.sellItem.Fee*sellOp.sellItem.DayExchangeRate,
			sellOp.sellItem.Fee*sellOp.sellItem.YearExchangeRate,
			report.Currency)
		report.TotalStockFee.Add(sellStockFee)
		report.TimeTestedStockRevenue.Add(sellOp.timeTestedRevenue)
		report.TimeTestedStockFee.Add(sellStockFee.MultiplyNew(timeTestedRatio))
	}

	for _, dividend := range dividends {
		report.DividendRevenue.Add(newAccountingValue(
			dividend.BrokerAmount*dividend.DayExchangeRate,
			dividend.BrokerAmount*dividend.YearExchangeRate, report.Currency))
		report.DividendFee.Add(newAccountingValue(
			dividend.Fee*dividend.DayExchangeRate,
			dividend.Fee*dividend.YearExchangeRate, report.Currency))
	}

	return &report
}

func calculateSellExpense(sellOp *SellOperation, availableBuyItems ItemToSellCollection) {

	timeTestDate := util.GetDateThreeYearsBefore(sellOp.sellItem.Date)

	quantityToBeSold := sellOp.sellItem.Quantity
	for _, itemToSell := range availableBuyItems {
		if itemToSell.availableQuantity <= 0.0 {
			continue
		}

		soldItem := &SoldItem{
			buyItem: itemToSell.buyItem,
			// 3 years time test
			timeTested: itemToSell.buyItem.Date.Sub(timeTestDate).Nanoseconds() < 0,
		}
		soldRatio := 0.0
		newAvailableQuantity := itemToSell.availableQuantity - quantityToBeSold
		if newAvailableQuantity >= 0.0 {
			// sell operation has all buys processed
			soldItem.soldQuantity = quantityToBeSold
			itemToSell.availableQuantity = newAvailableQuantity
			soldRatio = soldItem.soldQuantity / itemToSell.buyItem.Quantity
		} else {
			// some buy items are still required to be sold by this sell operation
			quantityToBeSold -= itemToSell.availableQuantity
			soldItem.soldQuantity = itemToSell.availableQuantity
			// noting remains to be sold in the buy item
			itemToSell.availableQuantity = 0.0
			soldRatio = soldItem.soldQuantity / itemToSell.buyItem.Quantity
		}

		// calculate purchase for this item
		soldItem.fifoBuyPrice = newAccountingValue(
			soldRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.DayExchangeRate,
			soldRatio*itemToSell.buyItem.BankAmount*itemToSell.buyItem.YearExchangeRate,
			nil)
		soldItem.fifoBuyFee = newAccountingValue(
			soldRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.DayExchangeRate,
			soldRatio*itemToSell.buyItem.Fee*itemToSell.buyItem.YearExchangeRate,
			nil)
		// calculate revenue for this buy item
		revenue := newAccountingValue(
			soldRatio*sellOp.sellItem.BrokerAmount*sellOp.sellItem.DayExchangeRate,
			soldRatio*sellOp.sellItem.BrokerAmount*sellOp.sellItem.YearExchangeRate,
			nil)
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

func convertToItemsToSell(purchases ingest.TransactionLogItems) (resItems ItemToSellCollection) {
	for _, buyItem := range purchases {
		resItems = append(resItems, &ItemToSell{
			buyItem:           buyItem,
			availableQuantity: buyItem.Quantity,
			soldByItems:       ingest.TransactionLogItems{},
		})
	}
	return
}

func convertToSellOperations(sales ingest.TransactionLogItems) (resItems SellOperationCollection) {
	for _, sellItem := range sales {
		resItems = append(resItems, &SellOperation{
			sellItem:          sellItem,
			soldItems:         SoldItemCollection{},
			timeTestedRevenue: newAccountingValue(0, 0, nil),
			totalRevenue:      newAccountingValue(0, 0, nil),
		})
	}
	return
}

func sortByDate(input *ingest.TransactionLog) {
	sort.Sort(ByDate(input.Sales))
	sort.Sort(ByDate(input.Purchases))
	sort.Sort(ByDate(input.Dividends))
}

func getSalesInYear(sellOperations SellOperationCollection, from time.Time, to time.Time) (ret SellOperationCollection) {
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

func getDividendsInYear(dividends ingest.TransactionLogItems, from time.Time, to time.Time) (ret ingest.TransactionLogItems) {
	fromExclusive := from.Add(-1 * time.Second)
	toExclusive := to.Add(1 * time.Second)
	isTransactionTimestampBetween := func(item *ingest.TransactionLogItem) bool {
		return item.Date.After(fromExclusive) && item.Date.Before(toExclusive)
	}
	return filterDividends(dividends, isTransactionTimestampBetween)
}

func filterSellOperations(list SellOperationCollection, test func(*SellOperation) bool) (ret SellOperationCollection) {
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

func filterDividends(list ingest.TransactionLogItems, test func(*ingest.TransactionLogItem) bool) (ret ingest.TransactionLogItems) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}
