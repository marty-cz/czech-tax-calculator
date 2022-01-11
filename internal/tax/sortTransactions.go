package tax

import (
	"sort"

	"github.com/marty-cz/czech-tax-calculator/internal/ingest"
)

// ByDate implements sort.Interface based on the Date field.
type ByDate ingest.TransactionLogItems

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortByDate(input *ingest.TransactionLog) {
	sort.Sort(ByDate(input.Sales))
	sort.Sort(ByDate(input.Purchases))
	sort.Sort(ByDate(input.Dividends))
	sort.Sort(ByDate(input.AdditionalIncomes))
	sort.Sort(ByDate(input.AdditionalFees))
}
