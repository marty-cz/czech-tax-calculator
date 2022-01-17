package tax

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
)

type SoldItem struct {
	BuyItem      *ingest.TransactionLogItem
	SoldQuantity float64
	TimeTested   bool
	FifoBuy      *ValueAndFee
	Revenue      *ValueAndFee
}

func (x *SoldItem) String() string {
	return fmt.Sprintf("buyItem:%+v soldQuantity:%v timeTested:%v fifoBuy:(%v)",
		x.BuyItem, x.SoldQuantity, x.TimeTested, x.FifoBuy)
}

type SoldItems []*SoldItem

type SellOperation struct {
	SellItem          *ingest.TransactionLogItem
	SoldItems         SoldItems
	timeTestedRevenue *AccountingValue
	totalRevenue      *AccountingValue
}

func (x *SellOperation) String() string {
	return fmt.Sprintf("sellItem:%+v totalRevenue:(%v) timeTestedRevenue:(%v) soldItems:[%+v]",
		x.SellItem, x.totalRevenue, x.timeTestedRevenue, &x.SoldItems)
}

type SellOperations []*SellOperation

func convertToSellOperations(sales ingest.TransactionLogItems) (resItems SellOperations) {
	for _, sellItem := range sales {
		resItems = append(resItems, &SellOperation{
			SellItem:          sellItem,
			SoldItems:         SoldItems{},
			timeTestedRevenue: newAccountingValue(0, 0, DEFAULT_CURRENCY),
			totalRevenue:      newAccountingValue(0, 0, DEFAULT_CURRENCY),
		})
	}
	return
}
