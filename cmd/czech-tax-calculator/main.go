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
	log.SetLevel(log.DebugLevel)
}

func main() {
	stockInputPath := flag.String("stock-input", "", "File path to input file with Stocks transaction records")
	cryptoInputPath := flag.String("crypto-input", "", "File path to input file with Crypto-currencies transaction records")
	flag.Parse()

	// pre-check of Year change rate to CZK availability
	for year := 2011; year <= time.Now().Year(); year++ {
		if val, err := util.GetCzkExchangeRateInYear(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), *util.USD); err != nil || val <= 0.0 {
			log.Warnf("missing or invalid an Year exchange rate for year '%d' - result for that year will not be accurate. Please fill it in exchangeRate.go", year)
		}
	}

	var stockTaxReports tax.Reports
	if *stockInputPath != "" {
		stockTransactions, err := ingest.ProcessStocks(*stockInputPath)
		if err != nil {
			log.Errorf("Cannot ingest stock input file '%s' due to: %s", *stockInputPath, err)
		} else {
			log.Infof("Stocks: Ingested")

			stockTaxReports, err = tax.Calculate(stockTransactions, "2021", true)
			if err != nil {
				log.Errorf("Cannot create stock tax report due to: %s", *stockInputPath, err)
			} else {
				log.Infof("Stocks tax: Calculated")
			}
		}
	}

	var cryptoTaxReports tax.Reports
	if *cryptoInputPath != "" {
		cryptoTransactions, err := ingest.ProcessCryptos(*cryptoInputPath)
		if err != nil {
			log.Errorf("Cannot ingest crypto input file '%s' due to: %s", *cryptoInputPath, err)
		} else {
			log.Infof("Cryptos: Ingested")

			cryptoTaxReports, err = tax.Calculate(cryptoTransactions, "2021", false)
			if err != nil {
				log.Errorf("Cannot create crypto tax report due to: %s", *cryptoInputPath, err)
			} else {
				log.Infof("Cryptos tax: Calculated")
				log.Debugf("%+v", cryptoTaxReports)

			}
		}
	}

	statements := make(map[int]*export.Statement)
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

	for _, statement := range statements {
		if err := export.ExportToExcel(statement, fmt.Sprintf("./tax-statement-%d.xlsx", statement.Year)); err != nil {
			log.Errorf("Cannot create excel statement for year '%d'", statement.Year)
		} else {
			log.Infof("Report: Created for '%d'", statement.Year)
		}
	}

}
