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

	var coordsEqSumDRDs string
	var coordsEqSumDRYs string
	var coordsEqSumDFDs string
	var coordsEqSumDFYs string
	var coordsEqSumDPDs string
	var coordsEqSumDPYs string
	i := 0

	row += 2
	w.WriteCell(sheet, row, col, "Dividends (details)")
	w.WriteCell(sheet, row, col+1, "Country")
	w.WriteCell(sheet, row, col+2, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+3, "with YEAR exchange rate")
	for country, dividendReport := range report.DividendReports {
		row++
		w.WriteCell(sheet, row, col, "Revenue")
		w.WriteCell(sheet, row, col+1, country)
		coordsDRD := w.WriteAccountingCell(sheet, row, col+2, dividendReport.RawRevenue.Value.ValueWithDayExchangeRate, dividendReport.RawRevenue.Value.Currency)
		coordsEqSumDRDs += "+" + coordsDRD
		coordsDRY := w.WriteAccountingCell(sheet, row, col+3, dividendReport.RawRevenue.Value.ValueWithYearExchangeRate, dividendReport.RawRevenue.Value.Currency)
		coordsEqSumDRYs += "+" + coordsDRY
		row++
		w.WriteCell(sheet, row, col, "Paid Tax")
		coordsDTD := w.WriteAccountingCell(sheet, row, col+2, dividendReport.PaidTax.ValueWithDayExchangeRate, dividendReport.PaidTax.Currency)
		coordsDTY := w.WriteAccountingCell(sheet, row, col+3, dividendReport.PaidTax.ValueWithYearExchangeRate, dividendReport.PaidTax.Currency)
		row++
		w.WriteCell(sheet, row, col, "Fees")
		coordsDFD := w.WriteAccountingCell(sheet, row, col+2, dividendReport.RawRevenue.Fee.ValueWithDayExchangeRate, dividendReport.RawRevenue.Fee.Currency)
		coordsEqSumDFDs += "+" + coordsDFD
		coordsDFY := w.WriteAccountingCell(sheet, row, col+3, dividendReport.RawRevenue.Fee.ValueWithYearExchangeRate, dividendReport.RawRevenue.Fee.Currency)
		coordsEqSumDFYs += "+" + coordsDFY
		row++
		w.WriteCell(sheet, row, col, "Profit")
		coordsDPD := w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s-%s", coordsDRD, coordsDTD, coordsDFD), report.Currency)
		coordsEqSumDPDs += "+" + coordsDPD
		coordsDPY := w.WriteAccountingEqCell(sheet, row, col+3, fmt.Sprintf("%s-%s-%s", coordsDRY, coordsDTY, coordsDFY), report.Currency)
		coordsEqSumDPYs += "+" + coordsDPY
		i++
	}

	row += 2
	w.WriteCell(sheet, row, col, "Dividends (summary)")
	w.WriteCell(sheet, row, col+1, "with DAY exchange rate")
	w.WriteCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.WriteCell(sheet, row, col, "Revenue")
	coordsDRD := w.WriteAccountingEqCell(sheet, row, col+1, coordsEqSumDRDs, report.Currency)
	coordsDRY := w.WriteAccountingEqCell(sheet, row, col+2, coordsEqSumDRYs, report.Currency)
	row++
	w.WriteCell(sheet, row, col, "Fees")
	w.WriteAccountingEqCell(sheet, row, col+1, coordsEqSumDFDs, report.Currency)
	w.WriteAccountingEqCell(sheet, row, col+2, coordsEqSumDFYs, report.Currency)
	row++
	w.WriteCell(sheet, row, col, "Profit")
	coordsDPD := w.WriteAccountingEqCell(sheet, row, col+1, coordsEqSumDPDs, report.Currency)
	coordsDPY := w.WriteAccountingEqCell(sheet, row, col+2, coordsEqSumDPYs, report.Currency)

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
	w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s+%s", coordsSRD, coordsTSRD, coordsDRD, coordsARD), report.Currency)
	w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s+%s", coordsSRY, coordsTSRY, coordsDRY, coordsARY), report.Currency)
	row++
	w.WriteCell(sheet, row, col, "Total Profit")
	// profit - Time Tested profit + Dividend Profit + Dividend (to pay tax) profit + Additional profit
	w.WriteAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s+%s", coordsSPD, coordsTSPD, coordsDPD, coordsAPD), report.Currency)
	w.WriteAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s+%s", coordsSPY, coordsTSPY, coordsDPY, coordsAPY), report.Currency)

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
