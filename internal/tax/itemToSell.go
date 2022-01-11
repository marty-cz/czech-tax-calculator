package tax

import (
	"fmt"

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

type ItemToSellCollection []*ItemToSell

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
