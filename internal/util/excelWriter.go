package util

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	File                   *excelize.File
	floatCellStyle         int
	dateCellStyle          int
	accountingCellStyleMap map[*Currency]int
}

func NewExcelWriter() *ExcelWriter {
	w := ExcelWriter{
		File:                   excelize.NewFile(),
		accountingCellStyleMap: make(map[*Currency]int),
	}

	w.floatCellStyle, _ = w.File.NewStyle(&excelize.Style{
		DecimalPlaces: 2,
	})
	dateFormat := "dd.mm.yyyy"
	w.dateCellStyle, _ = w.File.NewStyle(&excelize.Style{
		CustomNumFmt: &dateFormat,
	})
	for _, c := range SupportedCurrencies {
		currencyFormat := fmt.Sprintf("\"%s\" #,##0.00", c.Symbol)
		w.accountingCellStyleMap[c], _ = w.File.NewStyle(&excelize.Style{
			CustomNumFmt: &currencyFormat,
		})
	}
	return &w
}

func (w ExcelWriter) WriteCell(sheet string, row, col int, value interface{}) (coords string) {
	coords = GetExcelCoords(row, col)
	w.File.SetCellValue(sheet, coords, value)
	return
}

func (w ExcelWriter) WriteDateCell(sheet string, row, col int, date time.Time) (coords string) {
	coords = GetExcelCoords(row, col)
	w.File.SetCellValue(sheet, coords, date)
	w.File.SetCellStyle(sheet, coords, coords, w.dateCellStyle)
	return
}

func (w ExcelWriter) WriteFloatNumberCell(sheet string, row, col int, value float64) (coords string) {
	coords = GetExcelCoords(row, col)
	w.File.SetCellValue(sheet, coords, value)
	w.File.SetCellStyle(sheet, coords, coords, w.floatCellStyle)
	return
}

func (w ExcelWriter) WriteAccountingCell(sheet string, row, col int, value float64, currency *Currency) (coords string) {
	coords = GetExcelCoords(row, col)
	w.File.SetCellValue(sheet, coords, value)
	w.File.SetCellStyle(sheet, coords, coords, w.accountingCellStyleMap[currency])
	return
}

func (w ExcelWriter) WriteAccountingEqCell(sheet string, row, col int, equation string, currency *Currency) (coords string) {
	coords = GetExcelCoords(row, col)
	w.File.SetCellFormula(sheet, coords, equation)
	w.File.SetCellStyle(sheet, coords, coords, w.accountingCellStyleMap[currency])
	return
}
