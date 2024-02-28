package util

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// https://www.kodap.cz/cs/pro-vas/prehledy/jednotny-kurz/jednotne-kurzy-men-stanovene-ministerstvem-financi-prehled.html
// formatted into same format as CNB files https://www.cnb.cz/cs/casto-kladene-dotazy/Kurzy-devizoveho-trhu-na-www-strankach-CNB/
const MFCR_CZK_EXCHANGE_RATE_IN_YEARS string = `Země|Měna|Množství|Kód|2011|2012|2013|2014|2015|2016|2017|2018|2019|2020|2021|2022|2023|2024
Austrálie|dolar|1|AUD|18,33|20,19|18,75|18,73|18,41|18,22|17,81|16,21|15,96|16,00|16,26|16,21|14,67|0
Brazílie|real|1|BRL|10,53|9,96|9,03|8,87|7,40|7,13|7,25|5,94|5,81|4,46|4,03|4,55|4,45|0
Bulharsko|lev|1|BGN|12,58|12,84|13,31|14,09|13,94|13,83|13,44|13,13|13,12|13,55|13,11|12,55|12,23|0
Čína|remminbi|1|CNY|2,73|3,09|3,19|3,39|3,93|3,69|3,44|3,29|3,32|3,36|3,37|3,47|3,12|0
Dánsko|koruna|1|DKK|3,30|3,37|3,49|3,70|3,66|3,63|3,53|3,44|3,44|3,55|3,45|3,30|3,22|0
EMU|euro|1|EUR|24,60|25,12|26,03|27,55|27,27|27,04|26,29|25,68|25,66|26,50|25,65|24,54|23,97|0
Estonsko|koruna|1|EEK|-|-|-|-|-|-|-|-|-|-|-|-|-|-
Filipíny|peso|100|PHP|40,69|46,23|45,92|46,98|54,18|51,48|46,02|41,33|44,47|46,70|43,99|42,92|39,82|0
Hongkong|dolar|1|HKD|2,26|2,51|2,52|2,69|3,19|3,16|2,98|2,78|2,93|2,98|2,79|2,99|2,83|0
Chorvatsko|kuna|1|HRK|3,30|3,34|3,43|3,61|3,58|3,59|3,53|3,46|3,46|3,51|3,41|3,25|-|-
Indie|rupie|100|INR|37,57|36,44|33,41|34,18|38,43|36,46|35,69|31,82|32,58|31,24|29,38|29,69|26,80|0
Indonesie|rupie|1000|IDR|2,01|2,07|1,87|1,76|1,84|1,84|1,73|1,53|1,62|1,59|1,52|1,57|1,46|0
Island|koruna|100|ISK|-|-|-|-|-|-|-|20,24|18,75|17,11|17,18|17,31|16,16|0
Izrael|šekel|1|ILS|4,91|5,06|5,43|5,81|6,35|6,40|6,48|6,04|6,46|6,76|6,73|6,93|5,97|0
Japonsko|jen|100|JPY|22,18|24,34|20,03|19,62|20,41|22,50|20,71|19,75|21,05|21,76|19,69|17,79|15,67|0
Jihoafrická rep.|rand|1|ZAR|2,43|2,38|2,01|1,92|1,92|1,68|1,75|1,65|1,59|1,40|1,46|1,43|1,20|0
Jižní Korea|won|100|KRW|1,59|1,74|1,79|1,98|2,18|2,11|2,07|1,98|1,96|1,96|1,89|1,81|1,69|0
Kanada|dolar|1|CAD|17,83|19,48|18,91|18,85|19,16|18,54|17,87|16,74|17,32|17,23|17,33|17,93|16,40|0
Litva|litas|1|LTL|7,12|7,27|7,54|7,98|-|-|-|-|-|-|-|-|-|-
Lotyšsko|lat|1|LVL|34,84|36,01|37,08|-|-|-|-|-|-|-|-|-|-|-
Maďarsko|forint|100|HUF|8,79|8,72|8,74|8,89|8,81|8,67|8,50|8,03|7,88|7,49|7,15|6,26|6,30|0
Malajsie|ringgit|1|MYR|5,76|6,32|6,17|6,37|6,32|5,92|5,41|5,35|5,54|5,51|5,25|5,31|4,85|0
Mexiko|peso|1|MXN|1,41|1,48|1,52|1,56|1,55|1,31|1,23|1,13|1,19|1,08|1,07|1,17|1,26|0
MMF|SDR|1|XDR|27,82|29,83|29,73|31,64|34,49|34,01|32,20|30,81|31,67|32,25|30,91|31,21|29,54|0
Norsko|koruna|1|NOK|3,16|3,37|3,31|3,28|3,04|2,92|2,81|2,67|2,61|2,46|2,52|2,43|2,09|0
Nový Zéland|dolar|1|NZD|14,04|15,82|15,98|17,22|17,14|17,11|16,51|15,02|15,13|15,07|15,32|14,79|13,57|0
Polsko|zlotý|1|PLN|5,96|6,02|6,18|6,57|6,52|6,18|6,20|6,02|5,97|5,93|5,61|5,24|5,31|0
Rumunsko|nové leu|1|RON|5,80|5,64|5,90|6,21|6,14|6,02|5,75|5,51|5,40|5,47|5,21|4,97|4,84|0
Rusko|rubl|100|RUB|59,96|62,67|61,13|53,87|40,14|37,07|39,86|34,65|35,54|31,67|29,40|-|-|-
Singapur|dolar|1|SGD|14,03|15,63|15,61|16,44|17,92|17,74|16,86|16,14|16,82|16,79|16,17|16,97|16,49|0
Švédsko|koruna|1|SEK|2,73|2,89|3,00|3,02|2,92|2,86|2,73|2,50|2,42|2,53|2,53|2,30|2,09|0
Švýcarsko|frank|1|CHF|20,00|20,86|21,18|22,72|25,62|24,79|23,60|22,30|23,10|24,74|23,76|24,51|24,69|0
Thajsko|baht|100|THB|57,54|62,72|63,47|64,22|71,91|69,60|68,59|67,44|74,25|73,85|67,75|66,68|63,56|0
Turecko|lira|1|TRY|10,46|10,86|10,19|9,53|9,01|8,11|6,38|4,62|4,04|3,29|2,44|1,42|0,94|0
USA|dolar|1|USD|17,60|19,45|19,56|20,90|24,69|24,53|23,18|21,78|22,93|23,14|21,72|23,41|22,14|0
Velká Británie|libra|1|GBP|28,25|30,96|30,63|34,32|37,66|32,96|30,04|28,98|29,31|29,80|29,88|28,72|27,59|0
`

type Currency struct {
	Name   string
	Symbol string
}

func (c Currency) String() string {
	return c.Name
}

type ExchangeRate struct {
	FromCurrency  *Currency
	ToCurrency    *Currency
	RateInYearMap map[string]float64
}

var (
	EUR                 *Currency   = &Currency{Name: "EUR", Symbol: "€"}
	USD                 *Currency   = &Currency{Name: "USD", Symbol: "$"}
	CZK                 *Currency   = &Currency{Name: "CZK", Symbol: "Kč"}
	SupportedCurrencies []*Currency = []*Currency{EUR, USD, CZK}
)

func GetCurrencyByName(name string) (*Currency, error) {
	aName := strings.TrimSpace(name)
	if strings.EqualFold(EUR.Name, aName) {
		return EUR, nil
	}
	if strings.EqualFold(USD.Name, aName) {
		return USD, nil
	}
	if strings.EqualFold(CZK.Name, aName) {
		return CZK, nil
	}
	return nil, fmt.Errorf("unsupported currency '%s'", aName)
}

func GetCzkExchangeRateInYear(date time.Time, currency Currency) (float64, error) {
	if currency.Name == CZK.Name {
		return 1.0, nil
	}

	rate := -1.0
	yearString := date.Format("2006") // YYYY

	scanner := bufio.NewScanner(strings.NewReader(MFCR_CZK_EXCHANGE_RATE_IN_YEARS))
	firstLine := true
	exchangeIndex := -1
	for scanner.Scan() {
		txt := scanner.Text()
		if firstLine {
			exchangeIndex = indexOf(yearString, strings.Split(txt, "|"))
			if exchangeIndex < 0 {
				return rate, fmt.Errorf("no year %v in table header '%v'", yearString, txt)
			}
			firstLine = false
		} else {
			var err error
			if strings.Contains(txt, fmt.Sprintf("|%s|", currency.Name)) {
				if rate, err = getCzkRateFromCnbString(txt, exchangeIndex); err != nil {
					return rate, err
				}
				break
			}
		}
	}

	if rate < 0.0 {
		return rate, fmt.Errorf("exchange rate for currency '%v' not found", currency)
	}
	return rate, nil
}

// https://www.cnb.cz/cs/casto-kladene-dotazy/Kurzy-devizoveho-trhu-na-www-strankach-CNB/
const DATE_FORMAT_FOR_CNB_DAY string = "02.01.2006" // DD.MM.YYYY
const CNB_DAY_URL string = "https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date="

func GetCzkExchangeRateInDay(date time.Time, currency Currency) (float64, error) {
	if currency.Name == CZK.Name {
		return 1.0, nil
	}

	rate := -1.0
	dateString := date.Format(DATE_FORMAT_FOR_CNB_DAY)
	dateBeforeString := ""
	resp, err := http.Get(CNB_DAY_URL + dateString)
	if err != nil {
		return rate, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	firstLine := true
	for scanner.Scan() {
		txt := scanner.Text()
		if firstLine {
			if !strings.Contains(txt, dateString) {
				dateToCheck := date
				dateBeforeString = date.Format(DATE_FORMAT_FOR_CNB_DAY)
				for !isBusinessDayInCzechia(dateToCheck) {
					dateToCheck = dateToCheck.Add(-24 * time.Hour)
					dateBeforeString = dateToCheck.Format(DATE_FORMAT_FOR_CNB_DAY)
				}
				if !strings.Contains(txt, dateBeforeString) {
					return rate, fmt.Errorf("received response is not from day %v (or %v in case of non business day) but '%v'", dateString, dateBeforeString, txt)
				}
			}
			firstLine = false
		} else {
			if strings.Contains(txt, fmt.Sprintf("|%s|", currency.Name)) {
				if rate, err = getCzkRateFromCnbString(txt, -1); err != nil {
					return rate, err
				}
				break
			}
		}
	}

	if rate <= 0.0 {
		return rate, fmt.Errorf("exchange rate for currency '%v' not found", currency)
	}
	return rate, nil
}

func isBusinessDayInCzechia(date time.Time) bool {
	switch date.Weekday() {
	case time.Saturday, time.Sunday:
		return false
	default:
		return !isPublicHolidayInCzechia(date)
	}
}

// https://github.com/vaniocz/svatky-vanio-cz
const DATE_FORMAT_FOR_PUBLIC_HOLIDAY string = "2006-01-02" // YYYY-MM-DD
const CZECH_PUBLIC_HOLIDAY_URL string = "https://svatky.vanio.cz/api/"

func isPublicHolidayInCzechia(date time.Time) bool {
	dateString := date.Format(DATE_FORMAT_FOR_PUBLIC_HOLIDAY)
	resp, err := http.Get(CZECH_PUBLIC_HOLIDAY_URL + dateString)
	if err != nil {
		log.Warnf("cannot retrieve public holiday for Czechia on %v. Will be treat as general day. %s", date, err)
		return false
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	return strings.Contains(body, "isPublicHoliday: 1")
}

// Maďarsko|forint|100|HUF|8,79|8,72|8,74|8,89|8,81|8,67|8,50|8,03|7,88|7,49|7,15
// means that 100 CZK = 8,79 or 8,72 or ... HUF
func getCzkRateFromCnbString(cnbLine string, positionOfRateColumn int) (float64, error) {
	rate := -1.0
	splitLine := strings.Split(cnbLine, "|")
	if len(splitLine) < 5 {
		return rate, fmt.Errorf("unsupported format '%s' (expects 'COUNTRY|CURRENCY_HUMAN_NAME|MULTIPLICATOR|CURRENCY|RATE(s)')", cnbLine)
	}

	multiplicator, err := strconv.ParseFloat(strings.Replace(splitLine[2], ",", ".", 1), 64)
	if err != nil || multiplicator < 0.0 {
		return rate, fmt.Errorf("invalid multiplicator number '%s' (expects positive float or int number)", splitLine[2])
	}

	rateColIdx := 4
	if positionOfRateColumn >= 0 {
		rateColIdx = positionOfRateColumn
	}
	rate, err = strconv.ParseFloat(strings.Replace(splitLine[rateColIdx], ",", ".", 1), 64)
	if err != nil || rate < 0.0 {
		return rate, fmt.Errorf("invalid exchange rate number '%s' (expects positive float or int number)", splitLine[rateColIdx])
	}
	rate = rate / multiplicator
	return rate, nil
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found
}
