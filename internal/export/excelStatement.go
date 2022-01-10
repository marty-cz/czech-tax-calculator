package export

import (
	"fmt"

	"github.com/marty-cz/czech-tax-calculator/internal/tax"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
	"github.com/xuri/excelize/v2"
)

type work struct {
	file                *excelize.File
	year                int
	numberCellStyle     int
	accountingCellStyle int
}

func newWork(report *tax.Report) *work {
	w := work{
		file: excelize.NewFile(),
		year: report.Year.Year(),
	}

	numFormat := fmt.Sprintf("\"%s\" #,##0.00", report.Currency.Symbol)
	w.accountingCellStyle, _ = w.file.NewStyle(&excelize.Style{
		CustomNumFmt: &numFormat,
	})
	w.numberCellStyle, _ = w.file.NewStyle(&excelize.Style{
		DecimalPlaces: 2,
	})
	return &w
}

func (w work) String() string {
	return fmt.Sprintf("year %d", w.year)
}

func (w work) writeOverviewStatement(report *tax.Report) error {
	sheet := "Overview"
	// Create a new sheet.
	w.file.NewSheet(sheet)
	// Set value of a cell.
	row, col := 3, 0
	row += 2
	w.writeCell(sheet, row, col, "Stocks")
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Revenue")
	coordsSRD := w.writeAccountingCell(sheet, row, col+1, report.TotalStockRevenue.ValueWithDayExchangeRate, *report.TotalStockRevenue.Currency)
	coordsSRY := w.writeAccountingCell(sheet, row, col+2, report.TotalStockRevenue.ValueWithYearExchangeRate, *report.TotalStockRevenue.Currency)
	row++
	w.writeCell(sheet, row, col, "Expense")
	coordsSED := w.writeAccountingCell(sheet, row, col+1, report.TotalStockExpense.ValueWithDayExchangeRate, *report.TotalStockExpense.Currency)
	coordsSEY := w.writeAccountingCell(sheet, row, col+2, report.TotalStockExpense.ValueWithYearExchangeRate, *report.TotalStockExpense.Currency)
	row++
	w.writeCell(sheet, row, col, "Fees")
	coordsSFD := w.writeAccountingCell(sheet, row, col+1, report.TotalStockFee.ValueWithDayExchangeRate, *report.TotalStockFee.Currency)
	coordsSFY := w.writeAccountingCell(sheet, row, col+2, report.TotalStockFee.ValueWithYearExchangeRate, *report.TotalStockFee.Currency)
	row++
	w.writeCell(sheet, row, col, "Profit")
	coordsSPD := w.writeAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s-%s", coordsSRD, coordsSED, coordsSFD), *report.Currency)
	coordsSPY := w.writeAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s-%s", coordsSRY, coordsSEY, coordsSFY), *report.Currency)

	row += 2
	w.writeCell(sheet, row, col, "Time tested Stocks (3 year test)")
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Revenue")
	coordsTSRD := w.writeAccountingCell(sheet, row, col+1, report.TimeTestedStockRevenue.ValueWithDayExchangeRate, *report.TimeTestedStockRevenue.Currency)
	coordsTSRY := w.writeAccountingCell(sheet, row, col+2, report.TimeTestedStockRevenue.ValueWithYearExchangeRate, *report.TimeTestedStockRevenue.Currency)
	row++
	w.writeCell(sheet, row, col, "Expense")
	coordsTSED := w.writeAccountingCell(sheet, row, col+1, report.TimeTestedStockExpense.ValueWithDayExchangeRate, *report.TimeTestedStockExpense.Currency)
	coordsTSEY := w.writeAccountingCell(sheet, row, col+2, report.TimeTestedStockExpense.ValueWithYearExchangeRate, *report.TimeTestedStockExpense.Currency)
	row++
	w.writeCell(sheet, row, col, "Fees")
	coordsTSFD := w.writeAccountingCell(sheet, row, col+1, report.TimeTestedStockFee.ValueWithDayExchangeRate, *report.TimeTestedStockFee.Currency)
	coordsTSFY := w.writeAccountingCell(sheet, row, col+2, report.TimeTestedStockFee.ValueWithYearExchangeRate, *report.TimeTestedStockFee.Currency)
	row++
	w.writeCell(sheet, row, col, "Profit")
	coordsTSPD := w.writeAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s-%s", coordsTSRD, coordsTSED, coordsTSFD), *report.Currency)
	coordsTSPY := w.writeAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s-%s", coordsTSRY, coordsTSEY, coordsTSFY), *report.Currency)

	row += 2
	w.writeCell(sheet, row, col, "Dividends")
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Revenue")
	coordsDRD := w.writeAccountingCell(sheet, row, col+1, report.DividendRevenue.ValueWithDayExchangeRate, *report.DividendRevenue.Currency)
	coordsDRY := w.writeAccountingCell(sheet, row, col+2, report.DividendRevenue.ValueWithYearExchangeRate, *report.DividendRevenue.Currency)
	row++
	w.writeCell(sheet, row, col, "Fees")
	coordsDFD := w.writeAccountingCell(sheet, row, col+1, report.DividendFee.ValueWithYearExchangeRate, *report.DividendFee.Currency)
	coordsDFY := w.writeAccountingCell(sheet, row, col+2, report.DividendFee.ValueWithYearExchangeRate, *report.DividendFee.Currency)
	row++
	w.writeCell(sheet, row, col, "Profit")
	coordsDPD := w.writeAccountingEqCell(sheet, row, col+1, fmt.Sprintf("%s-%s", coordsDRD, coordsDFD), *report.Currency)
	coordsDPY := w.writeAccountingEqCell(sheet, row, col+2, fmt.Sprintf("%s-%s", coordsDRY, coordsDFY), *report.Currency)

	row, col = 0, 0
	w.writeCell(sheet, row, col, "Year")
	w.writeCell(sheet, row, col+1, report.Year.Year())
	row++
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Total Revenue")
	// Stock revenue - Time Tested stock revenue + Dividend revenue
	w.writeAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s", coordsSRD, coordsTSRD, coordsDRD), *report.Currency)
	w.writeAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s", coordsSRY, coordsTSRY, coordsDRY), *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Total Profit")
	// Stock profit - Time Tested stock profit + Dividend Profit
	w.writeAccountingEqCell(sheet, row, col+1, fmt.Sprintf("(%s-%s)+%s", coordsSPD, coordsTSPD, coordsDPD), *report.Currency)
	w.writeAccountingEqCell(sheet, row, col+2, fmt.Sprintf("(%s-%s)+%s", coordsSPY, coordsTSPY, coordsDPY), *report.Currency)

	return nil
}

func (w work) writeCell(sheet string, row, col int, value interface{}) (coords string) {
	colLetter := util.GetColumnLetter(col)
	coords = fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellValue(sheet, coords, value)
	return
}

func (w work) writeNumberCell(sheet string, row, col int, value float64) (coords string) {
	colLetter := util.GetColumnLetter(col)
	coords = fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellValue(sheet, coords, value)
	w.file.SetCellStyle(sheet, coords, coords, w.numberCellStyle)
	return
}

func (w work) writeAccountingCell(sheet string, row, col int, value float64, currency util.Currency) (coords string) {
	colLetter := util.GetColumnLetter(col)
	coords = fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellValue(sheet, coords, value)
	w.file.SetCellStyle(sheet, coords, coords, w.accountingCellStyle)
	return
}

func (w work) writeAccountingEqCell(sheet string, row, col int, equation string, currency util.Currency) (coords string) {
	colLetter := util.GetColumnLetter(col)
	coords = fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellFormula(sheet, coords, equation)
	w.file.SetCellStyle(sheet, coords, coords, w.accountingCellStyle)
	return
}

func ExportToExcel(report *tax.Report, exportFilePath string) error {
	w := newWork(report)

	if err := w.writeOverviewStatement(report); err != nil {
		return fmt.Errorf("cannot write overview statement for year '%v': %v", w, err)
	}
	// Delete "Sheet1"
	w.file.DeleteSheet(w.file.GetSheetName(0))
	if err := w.file.SaveAs(exportFilePath); err != nil {
		return fmt.Errorf("cannot save excel file '%s': %v", exportFilePath, err)
	}
	return nil
}
