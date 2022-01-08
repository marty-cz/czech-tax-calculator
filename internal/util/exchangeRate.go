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
	EUR *Currency = &Currency{Name: "EUR", Symbol: "€"}
	USD *Currency = &Currency{Name: "USD", Symbol: "$"}
	CZK *Currency = &Currency{Name: "CZK", Symbol: "Kč"}
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

/*
// ratePerYearString supported format YEAR:RATE;YEAR:RATE;
// ie. call("USD", "2019:27.93;2020:23,14;2021:21,72")
func LoadExchangeRatesToCzk(fromCurrencyName string, ratesInYearsString string) (rate *ExchangeRate, err error) {
	fromCurrency, err := GetCurrencyByName(fromCurrencyName)
	if err != nil {
		return nil, fmt.Errorf("from currency name is invalid: %s", err)
	}
	rate = &ExchangeRate{FromCurrency: fromCurrency, ToCurrency: CZK, RateInYearMap: make(map[string]float64)}

	for _, rateInYearString := range strings.Split(ratesInYearsString, ";") {
		rateInYear := strings.Split(rateInYearString, ":")
		if len(rateInYear) != 2 {
			return nil, fmt.Errorf("invalid format of '%s' - expected 'YEAR:RATE'", rateInYearString)
		}
		year, err := strconv.Atoi(rateInYear[0])
		if err != nil || (year >= 1000 && year <= 9999) {
			return nil, fmt.Errorf("invalid date '%s' - expected value in <1000;9999>", rateInYear[0])
		}
		rateValue, err := strconv.ParseFloat(strings.Replace(rateInYear[1], ",", ".", 1), 64)
		if err != nil || rateValue < 0.0 {
			return nil, fmt.Errorf("invalid rate number '%s' - expected positive float number", rateInYear[1])
		}
		rate.RateInYearMap[strconv.Itoa(year)] = rateValue
	}

	if len(rate.RateInYearMap) == 0 {
		return nil, fmt.Errorf("exchange rates are empty")
	}

	return
}
*/

// https://www.cnb.cz/cs/casto-kladene-dotazy/Kurzy-devizoveho-trhu-na-www-strankach-CNB/
const DATE_DORMAT_FOR_CNB string = "02.01.2006" // DD.MM.YYYY
const CNB_URL string = "https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date="

func GetCzkExchangeRateInDay(date time.Time, currency Currency) (float64, error) {
	rate := -1.0
	dateString := date.Format(DATE_DORMAT_FOR_CNB)
	dateBeforeString := date.Add(-24 * time.Hour).Format(DATE_DORMAT_FOR_CNB)
	resp, err := http.Get(CNB_URL + dateString)
	if err != nil {
		return rate, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	firstLine := true
	for scanner.Scan() {
		txt := scanner.Text()
		if firstLine {
			if !strings.Contains(txt, dateString) && (isPublicHolidayInCzechia(date) && strings.Contains(txt, dateBeforeString)) {
				return rate, fmt.Errorf("received response is not from %v (or %v in case of public holiday) but '%v'", dateString, dateBeforeString, txt)
			}
			firstLine = false
		} else {
			// Maďarsko|forint|100|HUF|6,817
			// means that 100 CZK = 6,817 HUF
			if strings.Contains(txt, fmt.Sprintf("|%s|", currency.Name)) {
				splitLine := strings.Split(txt, "|")
				if len(splitLine) != 5 {
					return rate, fmt.Errorf("unsupported format '%s' (expects 'COUNTRY|CURRENCY_HUMAN_NAME|MULTIPLICATOR|CURRENCY|RATE')", txt)
				}

				multiplicator, err := strconv.ParseFloat(strings.Replace(splitLine[2], ",", ".", 1), 64)
				if err != nil || multiplicator < 0.0 {
					return rate, fmt.Errorf("invalid multiplicator number '%s' (expects positive float or int number)", splitLine[2])
				}
				rate, err = strconv.ParseFloat(strings.Replace(splitLine[4], ",", ".", 1), 64)
				if err != nil || rate < 0.0 {
					return rate, fmt.Errorf("invalid exchange rate number '%s' (expects positive float or int number)", splitLine[4])
				}
				rate = rate / multiplicator
				break
			}
		}
	}
	return rate, nil
}

// https://github.com/vaniocz/svatky-vanio-cz
const DATE_DORMAT_FOR_PUBLIC_HOLIDAY string = "2006-01-02" // YYYY-MM-DD
const CZECH_PUBLIC_HOLIDAY_URL string = "https://svatky.vanio.cz/api/"

func isPublicHolidayInCzechia(date time.Time) bool {
	dateString := date.Format(DATE_DORMAT_FOR_PUBLIC_HOLIDAY)
	resp, err := http.Get(CZECH_PUBLIC_HOLIDAY_URL + dateString)
	if err != nil {
		log.Warnf("Cannot retrieve public holiday for Czechia on %v. Will be treat as general day. %s", date, err)
		return false
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	return strings.Contains(body, "<dt>isPublicHoliday</dt><dd>1</dd>")
}
