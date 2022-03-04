package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/marty-cz/czech-tax-calculator/internal/export"
	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
	"github.com/marty-cz/czech-tax-calculator/internal/tax"
	"github.com/marty-cz/czech-tax-calculator/internal/util"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	stockInputPath := flag.String("stock-input", "", "File path to input file with Stocks transaction records")
	cryptoInputPath := flag.String("crypto-input", "", "File path to input file with Crypto-currencies transaction records")
	targetYear := flag.String("year", fmt.Sprint(time.Now().Year()-1), "Target year for taxes")
	flag.Parse()

	// pre-check of Year change rate to CZK availability
	for year := 2011; year <= time.Now().Year(); year++ {
		if val, err := util.GetCzkExchangeRateInYear(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), *util.USD); err != nil || val <= 0.0 {
			log.Warnf("missing or invalid an Year exchange rate for year '%d' - result for that year will not be accurate. Please fill it in exchangeRate.go", year)
		}
	}

	// process input files
	stockTaxReports := createTaxReport(*stockInputPath, *targetYear, ingest.StockItemType, ingest.ProcessStocks)
	cryptoTaxReports := createTaxReport(*cryptoInputPath, *targetYear, ingest.CryptoItemType, ingest.ProcessCryptos)

	// write to output file
	statements := createStatementMap(stockTaxReports, cryptoTaxReports)
	for _, statement := range statements {
		if err := export.ExportToExcel(statement, fmt.Sprintf("./tax-statement-%d.xlsx", statement.Year)); err != nil {
			log.Errorf("cannot create excel statement for year '%d'", statement.Year)
		} else {
			log.Infof("report: Created for '%d'", statement.Year)
		}
	}

}

func createTaxReport(sourceFilePath string, targetYear string, itemTypeString string, ingestFn func(string) (*ingest.TransactionLog, error)) (taxReports tax.Reports) {
	if sourceFilePath != "" {
		transactions, err := ingestFn(sourceFilePath)
		if err != nil {
			log.Errorf("%ss: cannot ingest input file '%s' due to: %s", itemTypeString, sourceFilePath, err)
		} else {
			log.Infof("%ss: all ingested", itemTypeString)

			taxReports, err = tax.Calculate(transactions, targetYear, true)
			if err != nil {
				log.Errorf("%ss: cannot create tax report due to: %s", itemTypeString, err)
			} else {
				log.Infof("%ss tax: Calculated (reports count: %d)", itemTypeString, len(taxReports))
			}
		}
	}
	return
}

func createStatementMap(stockTaxReports, cryptoTaxReports tax.Reports) (statements map[int]*export.Statement) {
	statements = make(map[int]*export.Statement)
	for _, stockReport := range stockTaxReports {
		year := stockReport.Year.Year()
		if statements[year] == nil {
			statements[year] = &export.Statement{
				StockReport:  stockReport,
				CryptoReport: nil,
				Year:         year,
			}
		} else {
			statements[year].StockReport = stockReport
			statements[year].Year = year
		}
	}
	for _, cryptoReport := range cryptoTaxReports {
		year := cryptoReport.Year.Year()
		if statements[year] == nil {
			statements[year] = &export.Statement{
				StockReport:  nil,
				CryptoReport: cryptoReport,
				Year:         year,
			}
		} else {
			statements[year].CryptoReport = cryptoReport
			statements[year].Year = year
		}
	}
	return
}
