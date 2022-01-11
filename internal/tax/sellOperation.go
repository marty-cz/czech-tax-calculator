package tax

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
)

type SoldItem struct {
	buyItem      *ingest.TransactionLogItem
	soldQuantity float64
	timeTested   bool
	fifoBuy      *ValueAndFee
}

func (x *SoldItem) String() string {
	return fmt.Sprintf("buyItem:%+v soldQuantity:%v timeTested:%v fifoBuy:(%v)",
		x.buyItem, x.soldQuantity, x.timeTested, x.fifoBuy)
}

type SoldItems []*SoldItem

type SellOperation struct {
	sellItem          *ingest.TransactionLogItem
	soldItems         SoldItems
	timeTestedRevenue *AccountingValue
	totalRevenue      *AccountingValue
}

func (x *SellOperation) String() string {
	return fmt.Sprintf("sellItem:%+v totalRevenue:(%v) timeTestedRevenue:(%v) soldItems:[%+v]",
		x.sellItem, x.totalRevenue, x.timeTestedRevenue, &x.soldItems)
}

type SellOperations []*SellOperation

func convertToSellOperations(sales ingest.TransactionLogItems) (resItems SellOperations) {
	for _, sellItem := range sales {
		resItems = append(resItems, &SellOperation{
			sellItem:          sellItem,
			soldItems:         SoldItems{},
			timeTestedRevenue: newAccountingValue(0, 0, DEFAULT_CURRENCY),
			totalRevenue:      newAccountingValue(0, 0, DEFAULT_CURRENCY),
		})
	}
	return
}
