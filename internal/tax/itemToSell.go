package tax

import (
	"fmt"
	"strings"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
)

type ItemToSell struct {
	buyItem           *ingest.TransactionLogItem
	availableQuantity float64
	soldByItems       ingest.TransactionLogItems
}

func (x *ItemToSell) String() string {
	return fmt.Sprintf("buyItem:%+v availableQuantity:%v soldByItems:%+v", x.buyItem, x.availableQuantity, &x.soldByItems)
}

type ItemsToSell []*ItemToSell

func convertToItemsToSell(purchases ingest.TransactionLogItems) (resItems ItemsToSell) {
	for _, buyItem := range purchases {
		resItems = append(resItems, &ItemToSell{
			buyItem:           buyItem,
			availableQuantity: buyItem.Quantity,
			soldByItems:       ingest.TransactionLogItems{},
		})
	}
	return
}

func getAvailableItemsToSell(itemsToSell ItemsToSell, sellTransaction *ingest.TransactionLogItem) (ret ItemsToSell) {
	test := func(itemToSell *ItemToSell) bool {
		return strings.EqualFold(itemToSell.buyItem.Name, sellTransaction.Name) && itemToSell.availableQuantity > 0.0
	}
	return filterItemsToSell(itemsToSell, test)
}

func filterItemsToSell(list ItemsToSell, test func(*ItemToSell) bool) (ret ItemsToSell) {
	for _, item := range list {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}
