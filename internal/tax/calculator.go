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

type ItemToSell struct {
	buyItem           *ingest.TransactionLogItem
	availableQuantity float64
	soldByItems       ingest.TransactionLogItems
}

func (item *ItemToSell) String() string {
	return fmt.Sprintf("buyItem:%+v availableQuantity:%v soldByItems:%+v", item.buyItem, item.availableQuantity, &item.soldByItems)
}

type ItemToSellCollection []*ItemToSell

type SellOperation struct {
	sellItem                         *ingest.TransactionLogItem
	soldItems                        ingest.TransactionLogItems
	fifoBuyPriceWithDayExchangeRate  float64
	fifoBuyFeeWithDayExchangeRate    float64
	fifoBuyPriceWithYearExchangeRate float64
	fifoBuyFeeWithYearExchangeRate   float64
	currency                         *util.Currency
}

func (item *SellOperation) String() string {
	return fmt.Sprintf("sellItem:%+v fifoBuyPriceDayExchange:%v fifoBuyFeeDayExchang:%v fifoBuyPriceYearExchange:%v fifoBuyFeeYearExchang:%v currency:%v soldItems:%+v",
		item.sellItem,
		item.fifoBuyPriceWithDayExchangeRate, item.fifoBuyFeeWithDayExchangeRate, item.fifoBuyPriceWithYearExchangeRate, item.fifoBuyFeeWithYearExchangeRate,
		item.currency, &item.soldItems)
}

type SellOperationCollection []*SellOperation

type Tax struct {
}

func Calculate(transactions *ingest.TransactionLog, year string) (tax *Tax, err error) {

	var (
		dateStart, dateEnd time.Time
	)

	year = strings.TrimSpace(year)
	layout := "02.01.2006 15:04:05"
	if dateStart, err = time.Parse(layout, "01.01."+year+" 00:00:00"); err != nil {
		return nil, fmt.Errorf("invalid year '%s'", year)
	}
	if dateEnd, err = time.Parse(layout, "31.12."+year+" 23:59:59"); err != nil {
		return nil, fmt.Errorf("invalid year '%s", year)
	}

	sortByDate(transactions)

	itemsToSell := convertToItemsToSell(transactions.Purchases)
	sellOperations := convertToSellOperations(transactions.Sales)

	inYearSellOperations := getSalesInYear(sellOperations, dateStart, dateEnd)
	log.Debugf("Sale transaction count for year '%s' is '%d'", year, len(inYearSellOperations))
	for _, sellOp := range inYearSellOperations {
		availableBuyItems := getAvailableItemsToSell(itemsToSell, sellOp.sellItem)
		log.Debugf("For '%s' : %v", sellOp.sellItem.Name, availableBuyItems)
		processSell(sellOp, availableBuyItems)
		log.Debugf("Processed '%+v'", sellOp)
	}

	return &Tax{}, nil
}

// TODO: Should be the calculation of fifo buy price/fee rather based on
// percentage of sold buy item quantity? Because we have available data for
// itemToSell.buyItem.BrokerAmount and itemToSell.buyItem.BankAmount
// TODO2: Should be the prices/fees calculated for local currency (CZK) instead?
func processSell(sellOp *SellOperation, availableBuyItems ItemToSellCollection) {
	buyPriceWithDayExchangeRate := 0.0
	buyFeeWithDayExchangeRate := 0.0
	buyPriceWithYearExchangeRate := 0.0
	buyFeeWithYearExchangeRate := 0.0
	quantityToBeSold := sellOp.sellItem.Quantity
	for _, itemToSell := range availableBuyItems {
		if itemToSell.availableQuantity <= 0.0 {
			continue
		}
		newAvailableQuantity := itemToSell.availableQuantity - quantityToBeSold
		if newAvailableQuantity >= 0.0 {
			// sell operation has all buys processed
			itemToSell.availableQuantity = newAvailableQuantity

			buyPriceWithDayExchangeRate += quantityToBeSold * itemToSell.buyItem.ItemPrice * itemToSell.buyItem.DayExchangeRate
			buyFeeWithDayExchangeRate += quantityToBeSold * itemToSell.buyItem.Fee * itemToSell.buyItem.DayExchangeRate
			sellOp.fifoBuyPriceWithDayExchangeRate = buyPriceWithDayExchangeRate / sellOp.sellItem.Quantity
			sellOp.fifoBuyFeeWithDayExchangeRate = buyFeeWithDayExchangeRate / sellOp.sellItem.Quantity

			buyPriceWithYearExchangeRate += quantityToBeSold * itemToSell.buyItem.ItemPrice * itemToSell.buyItem.YearExchangeRate
			buyFeeWithYearExchangeRate += quantityToBeSold * itemToSell.buyItem.Fee * itemToSell.buyItem.YearExchangeRate
			sellOp.fifoBuyPriceWithYearExchangeRate = buyPriceWithYearExchangeRate / sellOp.sellItem.Quantity
			sellOp.fifoBuyFeeWithYearExchangeRate = buyFeeWithYearExchangeRate / sellOp.sellItem.Quantity

			sellOp.soldItems = append(sellOp.soldItems, itemToSell.buyItem)
			itemToSell.soldByItems = append(itemToSell.soldByItems, sellOp.sellItem)
			return
		} else {
			// some buy item are still required to be sold by this sell operation
			quantityToBeSold -= itemToSell.availableQuantity
			itemToSell.availableQuantity = 0.0

			buyPriceWithDayExchangeRate += itemToSell.buyItem.Quantity * itemToSell.buyItem.ItemPrice * itemToSell.buyItem.DayExchangeRate
			buyFeeWithDayExchangeRate += itemToSell.buyItem.Quantity * itemToSell.buyItem.Fee * itemToSell.buyItem.DayExchangeRate

			buyPriceWithYearExchangeRate += itemToSell.buyItem.Quantity * itemToSell.buyItem.ItemPrice * itemToSell.buyItem.YearExchangeRate
			buyFeeWithYearExchangeRate += itemToSell.buyItem.Quantity * itemToSell.buyItem.Fee * itemToSell.buyItem.YearExchangeRate
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
			sellItem:  sellItem,
			soldItems: ingest.TransactionLogItems{},
			currency:  util.CZK,
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
