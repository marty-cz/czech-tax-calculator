package util

import (
	"fmt"
	"strings"
)

func IsRowEmpty(row []string, nameColIndex int) bool {
	if nameColIndex >= 0 && nameColIndex < len(row) {
		return strings.TrimSpace(row[nameColIndex]) == ""
	}
	return true
}

func ValidateTableHeader(row []string, legend map[string]int) (err error) {

	if len(row) != len(legend) {
		return fmt.Errorf("unexpected count of columns '%d', but should be '%d'", len(row), len(legend))
	}

	for expColName, expColNo := range legend {
		if err := throwErrorIfBadColumnName(row, expColNo, expColName); err != nil {
			return err
		}
	}

	return nil
}

func ConvertCurrency(rawCurencyStr string) (currency string, err error) {

	if strings.Contains(rawCurencyStr, "USD") {
		return "USD", nil
	}
	if strings.Contains(rawCurencyStr, "EUR") {
		return "EUR", nil
	}
	if strings.Contains(rawCurencyStr, "CZK") {
		return "CZK", nil
	}

	return "", fmt.Errorf("unrecognized currency '%s'", rawCurencyStr)
}

func throwErrorIfBadColumnName(row []string, index int, expected string) (err error) {
	if !strings.EqualFold(row[index], expected) {
		return fmt.Errorf("column name at %d has unexpected name '%s', but should be '%s'", index, row[index], expected)
	}
	return nil
}
