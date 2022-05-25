package ticker

import (
	"time"

	"gorm.io/gorm"
)

func AvgPrice(open, close, high, low float64) float64 {
	return (open + close + high + low) / 4
}

func AvgTime(openTime, endTime int64) int64 {
	return (openTime + endTime) / 2
}

type TickerPrice struct {
	OpenTime  int64   `json:"open_time"`
	CloseTime int64   `json:"close_time"`
	Symbol    string  `json:"symbol"`
	Source    string  `json:"source"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Interval  string  `json:"interval"`
	AvgTime   int64   `json:"avg_time"`
}

type Dataset struct {
	gorm.Model

	Set       string
	Name      string
	Value     float64
	Timestamp time.Time
	UnixTime  int64 `gorm:"index,unique"`
}

func BuildPriceDataset(tps []TickerPrice) []Dataset {
	var dss []Dataset
	for _, tp := range tps {
		dss = append(dss, Dataset{
			Set:       tp.Symbol,
			Name:      tp.Symbol + "-" + "price",
			Value:     tp.Price,
			Timestamp: time.UnixMilli(tp.AvgTime),
			UnixTime:  tp.AvgTime,
		})
	}
	return dss
}

func BuildVolDataSet(tps []TickerPrice) []Dataset {
	var dss []Dataset
	for _, tp := range tps {
		dss = append(dss, Dataset{
			Set:       tp.Symbol,
			Name:      tp.Symbol + "-" + "volume",
			Value:     tp.Volume,
			Timestamp: time.UnixMilli(tp.AvgTime),
			UnixTime:  tp.AvgTime,
		})
	}
	return dss
}
