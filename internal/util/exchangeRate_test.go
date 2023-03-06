package util

import (
	"testing"
	"time"
)

func createDate(day, month, year int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestGetCzkExchangeRateInDay(t *testing.T) {
	type args struct {
		date     time.Time
		currency Currency
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=27.12.2020
		{"Sunday 27.12. -> use 23.12.", args{date: createDate(27, 12, 2020), currency: *EUR}, 26.370, false},
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=09.01.2021
		{"Sunday 10.01. -> use 08.01.", args{createDate(10, 1, 2021), *EUR}, 26.165, false},
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=11.01.2021
		{"Monday 11.01. -> use 11.01.", args{createDate(11, 1, 2021), *EUR}, 26.240, false},
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=11.01.2021
		{"Not known currency", args{createDate(11, 1, 2021), Currency{Name: "XYZ", Symbol: "?"}}, -1.0, true},
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=11.01.1821
		{"Date out of range", args{createDate(11, 1, 1821), *EUR}, -1.0, true},
		// https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=11.01.1821
		{"CZK exchange rate is 1.0", args{createDate(11, 1, 1821), *CZK}, 1.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCzkExchangeRateInDay(tt.args.date, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCzkExchangeRateInDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCzkExchangeRateInDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
