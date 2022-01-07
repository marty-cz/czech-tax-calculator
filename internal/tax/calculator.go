package tax

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
	log "github.com/sirupsen/logrus"
)

// ByAge implements sort.Interface based on the Age field.
type ByDate []*ingest.TransactionLogItem

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ItemToSell struct {
	buyItem           *ingest.TransactionLogItem
	availableQuantity float64
	soldByItems       []*ingest.TransactionLogItem
}

type SellOperation struct {
	sellItem  *ingest.TransactionLogItem
	soldItems []*ingest.TransactionLogItem
}

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
		log.Debugf("For '%s' : %s", sellOp.sellItem.Name, availableBuyItems)
	}

	return &Tax{}, nil
}

func convertToItemsToSell(purchases []*ingest.TransactionLogItem) (resItems []*ItemToSell) {
	for _, buyItem := range purchases {
		resItems = append(resItems, &ItemToSell{
			buyItem:           buyItem,
			availableQuantity: buyItem.Quantity,
			soldByItems:       []*ingest.TransactionLogItem{},
		})
	}
	return
}

func convertToSellOperations(sales []*ingest.TransactionLogItem) (resItems []*SellOperation) {
	for _, sellItem := range sales {
		resItems = append(resItems, &SellOperation{
			sellItem:  sellItem,
			soldItems: []*ingest.TransactionLogItem{},
		})
	}
	return
}

func sortByDate(input *ingest.TransactionLog) {
	sort.Sort(ByDate(input.Sales))
	sort.Sort(ByDate(input.Purchases))
	sort.Sort(ByDate(input.Dividends))
}

func getSalesInYear(sellOperations []*SellOperation, from time.Time, to time.Time) (ret []*SellOperation) {
	fromExclusive := from.Add(-1 * time.Second)
	toExclusive := to.Add(1 * time.Second)
	isTransactionTimestampBetween := func(item *SellOperation) bool {
		return item.sellItem.Date.After(fromExclusive) && item.sellItem.Date.Before(toExclusive)
	}
	return filterSellOperations(sellOperations, isTransactionTimestampBetween)
}

func getAvailableItemsToSell(itemsToSell []*ItemToSell, sellTransaction *ingest.TransactionLogItem) (ret []*ItemToSell) {
	test := func(itemToSell *ItemToSell) bool {
		return strings.EqualFold(itemToSell.buyItem.Name, sellTransaction.Name) && itemToSell.availableQuantity > 0.0
	}
	return filterItemsToSell(itemsToSell, test)
}

func filterSellOperations(list []*SellOperation, test func(*SellOperation) bool) (ret []*SellOperation) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

func filterItemsToSell(list []*ItemToSell, test func(*ItemToSell) bool) (ret []*ItemToSell) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}
