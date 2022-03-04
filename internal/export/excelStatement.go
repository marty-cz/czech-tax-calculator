package export

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/tax"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

type Statement struct {
	StockReport  *tax.Report
	CryptoReport *tax.Report
	Year         int
}

func writeOverviewStatement(w *util.ExcelWriter, report *tax.Report, itemTypeString string) error {
	sheet := "Overview - " + itemTypeString
	// Create a new sheet.
	w.File.NewSheet(sheet)
	// Set value of a cell.
	row, col := 3, 0
	row += 2
	w.WriteCell(sheet, row, col, itemTypeString)
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsSRD := w.WriteAccountingCell(sheet, row, col+1, report.TotalItemRevenue.ValueWithDayExchangeRate, report.TotalItemRevenue.Currency)
	coordsSRY := w.WriteAccountingCell(sheet, row, col+2, report.TotalItemRevenue.ValueWithYearExchangeRate, report.TotalItemRevenue.Currency)
	row++
	w.WriteCell(sheet, row, col, "Expense")
	coordsSED := w.WriteAccountingCell(sheet, row, col+1, report.TotalItemFifoExpense.Value.ValueWithDayExchangeRate, report.TotalItemFifoExpense.Value.Currency)
	coordsSEY := w.WriteAccountingCell(sheet, row, col+2, report.TotalItemFifoExpense.Value.ValueWithYearExchangeRate, report.TotalItemFifoExpense.Value.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	coordsSFD := w.WriteAccountingCell(sheet, row, col+1, report.TotalItemFifoExpense.Fee.ValueWithDayExchangeRate, report.TotalItemFifoExpense.Fee.Currency)
	coordsSFY := w.WriteAccountingCell(sheet, row, col+2, report.TotalItemFifoExpense.Fee.ValueWithYearExchangeRate, report.TotalItemFifoExpense.Fee.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsSPD := w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s-%s", coordsSRD, coordsSED, coordsSFD), report.Currency)
	coordsSPY := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s-%s", coordsSRY, coordsSEY, coordsSFY), report.Currency)

	row += 2
	w.WriteCell(sheet, row, col, "Time tested "+itemTypeString+" (3 year test)")
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsTSRD := w.WriteAccountingCell(sheet, row, col+1, report.TimeTestedItemRevenue.ValueWithDayExchangeRate, report.TimeTestedItemRevenue.Currency)
	coordsTSRY := w.WriteAccountingCell(sheet, row, col+2, report.TimeTestedItemRevenue.ValueWithYearExchangeRate, report.TimeTestedItemRevenue.Currency)
	row++
	w.WriteCell(sheet, row, col, "Expense")
	coordsTSED := w.WriteAccountingCell(sheet, row, col+1, report.TimeTestedItemFifoExpense.Value.ValueWithDayExchangeRate, report.TimeTestedItemFifoExpense.Value.Currency)
	coordsTSEY := w.WriteAccountingCell(sheet, row, col+2, report.TimeTestedItemFifoExpense.Value.ValueWithYearExchangeRate, report.TimeTestedItemFifoExpense.Value.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	coordsTSFD := w.WriteAccountingCell(sheet, row, col+1, report.TimeTestedItemFifoExpense.Fee.ValueWithDayExchangeRate, report.TimeTestedItemFifoExpense.Fee.Currency)
	coordsTSFY := w.WriteAccountingCell(sheet, row, col+2, report.TimeTestedItemFifoExpense.Fee.ValueWithYearExchangeRate, report.TimeTestedItemFifoExpense.Fee.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsTSPD := w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s-%s", coordsTSRD, coordsTSED, coordsTSFD), report.Currency)
	coordsTSPY := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s-%s", coordsTSRY, coordsTSEY, coordsTSFY), report.Currency)

	row += 2
	w.WriteCell(sheet, row, col, "Dividends (tax payed already)")
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsDRD := w.WriteAccountingCell(sheet, row, col+1, report.DividendRevenue.Value.ValueWithDayExchangeRate, report.DividendRevenue.Value.Currency)
	coordsDRY := w.WriteAccountingCell(sheet, row, col+2, report.DividendRevenue.Value.ValueWithYearExchangeRate, report.DividendRevenue.Value.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	coordsDFD := w.WriteAccountingCell(sheet, row, col+1, report.DividendRevenue.Fee.ValueWithYearExchangeRate, report.DividendRevenue.Fee.Currency)
	coordsDFY := w.WriteAccountingCell(sheet, row, col+2, report.DividendRevenue.Fee.ValueWithYearExchangeRate, report.DividendRevenue.Fee.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsDPD := w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s", coordsDRD, coordsDFD), report.Currency)
	coordsDPY := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s", coordsDRY, coordsDFY), report.Currency)

	row += 2
	w.WriteCell(sheet, row, col, "Dividends (to pay tax from)")
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsDTRD := w.WriteAccountingCell(sheet, row, col+1, report.DividendRevenueWithNoTaxPayed.Value.ValueWithDayExchangeRate, report.DividendRevenueWithNoTaxPayed.Value.Currency)
	coordsDTRY := w.WriteAccountingCell(sheet, row, col+2, report.DividendRevenueWithNoTaxPayed.Value.ValueWithYearExchangeRate, report.DividendRevenueWithNoTaxPayed.Value.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	coordsDTFD := w.WriteAccountingCell(sheet, row, col+1, report.DividendRevenueWithNoTaxPayed.Fee.ValueWithYearExchangeRate, report.DividendRevenueWithNoTaxPayed.Fee.Currency)
	coordsDTFY := w.WriteAccountingCell(sheet, row, col+2, report.DividendRevenueWithNoTaxPayed.Fee.ValueWithYearExchangeRate, report.DividendRevenueWithNoTaxPayed.Fee.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsDTPD := w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s", coordsDTRD, coordsDTFD), report.Currency)
	coordsDTPY := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s", coordsDTRY, coordsDTFY), report.Currency)

	row += 2
	w.WriteCell(sheet, row, col, "Additional")
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsARD := w.WriteAccountingCell(sheet, row, col+1, report.AdditionalRevenue.Value.ValueWithDayExchangeRate, report.AdditionalRevenue.Value.Currency)
	coordsARY := w.WriteAccountingCell(sheet, row, col+2, report.AdditionalRevenue.Value.ValueWithYearExchangeRate, report.AdditionalRevenue.Value.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	coordsAFD := w.WriteAccountingCell(sheet, row, col+1, report.AdditionalRevenue.Fee.ValueWithYearExchangeRate, report.AdditionalRevenue.Fee.Currency)
	coordsAFY := w.WriteAccountingCell(sheet, row, col+2, report.AdditionalRevenue.Fee.ValueWithYearExchangeRate, report.AdditionalRevenue.Fee.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsAPD := w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s", coordsARD, coordsAFD), report.Currency)
	coordsAPY := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s", coordsARY, coordsAFY), report.Currency)

	row, col = 0, 0
	w.WriteCell(sheet, row, col, "Year")
	w.WriteCell(sheet, row, col+1, report.Year.Year())
	row++
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Total Revenue")
	// revenue - Time Tested revenue + Dividend revenue + Dividend (to pay tax) revenue + Additional revenue
	w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s+%s+%s", coordsSRD, coordsTSRD, coordsDRD, coordsDTRD, coordsARD), report.Currency)
	w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s+%s+%s", coordsSRY, coordsTSRY, coordsDRY, coordsDTRY, coordsARY), report.Currency)
	row++
	w.WriteCell(sheet, row, col, "Total Profit")
	// profit - Time Tested profit + Dividend Profit + Dividend (to pay tax) profit + Additional profit
	w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s+%s+%s", coordsSPD, coordsTSPD, coordsDPD, coordsDTPD, coordsAPD), report.Currency)
	w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s+%s+%s", coordsSPY, coordsTSPY, coordsDPY, coordsDTPY, coordsAPY), report.Currency)

	return nil
}

func ExportToExcel(statement *Statement, exportFilePath string) error {
	w := util.NewExcelWriter()

	// write overviews
	if statement.StockReport != nil {
		if err := writeOverviewStatement(w, statement.StockReport, "Stocks"); err != nil {
			return fmt.Errorf("cannot write stock overview statement for year '%v': %v", w, err)
		}
	}
	if statement.CryptoReport != nil {
		if err := writeOverviewStatement(w, statement.CryptoReport, "Cryptos"); err != nil {
			return fmt.Errorf("cannot write crypto overview statement for year '%v': %v", w, err)
		}
	}
	// write sales log
	if statement.StockReport != nil {
		if err := salesLogToExcel(w, statement.StockReport.SellOperations, "Stocks"); err != nil {
			return fmt.Errorf("cannot write stock sales log for year '%v': %v", w, err)
		}
	}
	if statement.CryptoReport != nil {
		if err := salesLogToExcel(w, statement.CryptoReport.SellOperations, "Cryptos"); err != nil {
			return fmt.Errorf("cannot write crypto sales log for year '%v': %v", w, err)
		}
	}
	// Delete "Sheet1"
	w.File.DeleteSheet(w.File.GetSheetName(0))
	if err := w.File.SaveAs(exportFilePath); err != nil {
		return fmt.Errorf("cannot save excel file '%s': %v", exportFilePath, err)
	}
	return nil
}
