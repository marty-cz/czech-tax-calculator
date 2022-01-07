package ingest

import (
	log "github.com/sirupsen/logrus"
	//	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
)

func ProcessCryptos(filePath string) (err error) {

	log.Debugf("Processing cryptos input file '%s'", filePath)
	log.Warn("Crypto currencies are not currently supported")
	return nil
}
