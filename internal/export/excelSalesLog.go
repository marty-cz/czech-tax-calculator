package export

import (
	"github.com/marty-cz/czech-tax-calculator/internal/tax"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

func salesLogToExcel(w *util.ExcelWriter, sales tax.SellOperations, itemTypeString string) error {

	sheet := "Sales - " + itemTypeString
	// Create a new sheet.
	w.File.NewSheet(sheet)

	// write header
	row, col := 0, 0
	w.WriteCell(sheet, row, col, itemTypeString)
	w.WriteCell(sheet, row, col+1, "Sell Date")
	w.WriteCell(sheet, row, col+2, "Sell Price (Day ExR, FIFO)")
	w.WriteCell(sheet, row, col+3, "Sell Price (Year ExR, FIFO)")
	w.WriteCell(sheet, row, col+4, "Buy Date")
	w.WriteCell(sheet, row, col+5, "Time Tested")
	w.WriteCell(sheet, row, col+6, "Buy Price (Day ExR, FIFO)")
	w.WriteCell(sheet, row, col+7, "Buy Price (Year ExR, FIFO)")
	w.WriteCell(sheet, row, col+8, "Fee (Day ExR, FIFO)")
	w.WriteCell(sheet, row, col+9, "Fee (Year ExR, FIFO)")

	// write log
	for _, sellOp := range sales {
		col := 0
		for _, soldItem := range sellOp.SoldItems {
			row++
			w.WriteCell(sheet, row, col, sellOp.SellItem.Name)
			w.WriteDateCell(sheet, row, col+1, sellOp.SellItem.Date)
			w.WriteAccountingCell(sheet, row, col+2, soldItem.Revenue.Value.ValueWithDayExchangeRate, soldItem.Revenue.Value.Currency)
			w.WriteAccountingCell(sheet, row, col+3, soldItem.Revenue.Value.ValueWithYearExchangeRate, soldItem.Revenue.Value.Currency)
			w.WriteDateCell(sheet, row, col+4, soldItem.BuyItem.Date)
			w.WriteCell(sheet, row, col+5, soldItem.TimeTested)
			w.WriteAccountingCell(sheet, row, col+6, soldItem.FifoBuy.Value.ValueWithDayExchangeRate, soldItem.FifoBuy.Value.Currency)
			w.WriteAccountingCell(sheet, row, col+7, soldItem.FifoBuy.Value.ValueWithYearExchangeRate, soldItem.FifoBuy.Value.Currency)
			w.WriteAccountingCell(sheet, row, col+8, soldItem.FifoBuy.Fee.ValueWithDayExchangeRate+soldItem.Revenue.Fee.ValueWithDayExchangeRate, soldItem.FifoBuy.Value.Currency)
			w.WriteAccountingCell(sheet, row, col+9, soldItem.FifoBuy.Fee.ValueWithYearExchangeRate+soldItem.Revenue.Fee.ValueWithYearExchangeRate, soldItem.FifoBuy.Value.Currency)
		}
	}

	return nil
}
