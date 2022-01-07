package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
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

	if *stockInputPath != "" {
		err := ingest.ProcessStocks(*stockInputPath)
		if err != nil {
			log.Errorf("Cannot ingest stock input file '%s' due to: %s", *stockInputPath, err)
		}
	}

	if *cryptoInputPath != "" {
		err := ingest.ProcessCryptos(*cryptoInputPath)
		if err != nil {
			log.Fatalf("Cannot ingest crypto input file '%s' due to: %s", *cryptoInputPath, err)
		}
	}

}
