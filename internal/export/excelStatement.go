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
	row, col := 0, 0

	w.writeCell(sheet, row, col, "Year")
	w.writeCell(sheet, row, col+1, report.Year.Year())
	row++
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Total Revenue")
	w.writeAccountingCell(sheet, row, col+1, report.DividendRevenueWithDayExchangeRate+report.RevenueWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.DividendRevenueWithYearExchangeRate+report.RevenueWithYearExchangeRate, *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Total Profit")
	w.writeAccountingCell(sheet, row, col+1, report.DividendRevenueWithDayExchangeRate-report.DividendFeeWithDayExchangeRate+report.RevenueWithDayExchangeRate-report.ExpenseWithDayExchangeRate-report.FeeWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.DividendRevenueWithYearExchangeRate-report.DividendFeeWithYearExchangeRate+report.RevenueWithYearExchangeRate-report.ExpenseWithYearExchangeRate-report.FeeWithYearExchangeRate, *report.Currency)

	row += 2
	w.writeCell(sheet, row, col, "Stocks")
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Revenue")
	w.writeAccountingCell(sheet, row, col+1, report.RevenueWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.RevenueWithYearExchangeRate, *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Expense")
	w.writeAccountingCell(sheet, row, col+1, report.ExpenseWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.ExpenseWithYearExchangeRate, *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Fees")
	w.writeAccountingCell(sheet, row, col+1, report.FeeWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.FeeWithYearExchangeRate, *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Profit")
	w.writeAccountingCell(sheet, row, col+1, report.RevenueWithDayExchangeRate-report.ExpenseWithDayExchangeRate-report.FeeWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.RevenueWithYearExchangeRate-report.ExpenseWithYearExchangeRate-report.FeeWithYearExchangeRate, *report.Currency)

	row += 2
	w.writeCell(sheet, row, col, "Dividends")
	w.writeCell(sheet, row, col+1, "with DAY exchange rate")
	w.writeCell(sheet, row, col+2, "with YEAR exchange rate")
	row++
	w.writeCell(sheet, row, col, "Revenue")
	w.writeAccountingCell(sheet, row, col+1, report.DividendRevenueWithDayExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.DividendRevenueWithYearExchangeRate, *report.Currency)
	row++
	w.writeCell(sheet, row, col, "Fees")
	w.writeAccountingCell(sheet, row, col+1, report.DividendFeeWithYearExchangeRate, *report.Currency)
	w.writeAccountingCell(sheet, row, col+2, report.DividendFeeWithYearExchangeRate, *report.Currency)

	return nil
}

func (w work) writeCell(sheet string, row, col int, value interface{}) {
	colLetter := util.GetColumnLetter(col)
	w.file.SetCellValue(sheet, fmt.Sprintf("%s%d", colLetter, row+1), value)
}

func (w work) writeNumberCell(sheet string, row, col int, value float64) {
	colLetter := util.GetColumnLetter(col)
	coords := fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellValue(sheet, coords, value)
	w.file.SetCellStyle(sheet, coords, coords, w.numberCellStyle)
}

func (w work) writeAccountingCell(sheet string, row, col int, value float64, currency util.Currency) {
	colLetter := util.GetColumnLetter(col)
	coords := fmt.Sprintf("%s%d", colLetter, row+1)
	w.file.SetCellValue(sheet, coords, value)
	w.file.SetCellStyle(sheet, coords, coords, w.accountingCellStyle)
}

func ExportToExcel(report *tax.Report, exportFilePath string) error {
	w := newWork(report)

	if err := w.writeOverviewStatement(report); err != nil {
		return fmt.Errorf("cannot write overview statement for year '%v': %v", w, err)
	}
	if err := w.file.SaveAs(exportFilePath); err != nil {
		return fmt.Errorf("cannot save excel file '%s': %v", exportFilePath, err)
	}
	return nil
}
