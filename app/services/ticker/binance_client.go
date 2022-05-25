package ticker

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"

	"github.com/namnm1991/hackamoon07/app/services/utils"
)

var (
	apiKey    = "your api key"
	secretKey = "your secret key"
	baseUrl   = "https://www.binance.com"
)

type BinanceClient struct {
	baseClient *binance.Client
	fapiUrl    string
}

func NewBinanceClient() *BinanceClient {
	c := binance.NewClient(apiKey, secretKey)
	return &BinanceClient{
		baseClient: c,
		fapiUrl:    fmt.Sprintf("%s/fapi/v1", baseUrl),
	}
}

func (b *BinanceClient) GetTickerPrice(symbol, interval string, startTime, endTime int64) ([]TickerPrice, error) {
	log.Printf("https://www.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d", symbol, interval, startTime, endTime)
	klines, err := b.baseClient.NewKlinesService().Symbol(symbol).StartTime(startTime).EndTime(endTime).Interval(interval).Do(context.Background())
	if err != nil {
		return []TickerPrice{}, err
	}
	results := make([]TickerPrice, len(klines))
	for i, kl := range klines {
		closePrice, err := strconv.ParseFloat(kl.Close, 64)
		if err != nil {
			continue
		}
		openPrice, err := strconv.ParseFloat(kl.Open, 64)
		if err != nil {
			continue
		}
		lowPrice, err := strconv.ParseFloat(kl.Low, 64)
		if err != nil {
			continue
		}
		highPrice, err := strconv.ParseFloat(kl.High, 64)
		if err != nil {
			continue
		}
		vol, err := strconv.ParseFloat(kl.Volume, 64)
		if err != nil {
			continue
		}
		results[i] = TickerPrice{
			Symbol:    strings.ToUpper(symbol),
			Volume:    vol,
			OpenTime:  kl.OpenTime,
			CloseTime: kl.CloseTime,
			Price:     AvgPrice(openPrice, closePrice, highPrice, lowPrice),
			Interval:  interval,
			Source:    "binance-exchange",
			AvgTime:   AvgTime(kl.OpenTime, kl.CloseTime),
		}
	}
	return results, nil
}

func (b *BinanceClient) GetListFundingRate(symbol string, startTime int64) ([]TickerFundingRate, error) {
	url := fmt.Sprintf("%s/fundingRate", b.fapiUrl)
	var resp []TickerFundingRate
	params := map[string]string{
		"symbol":    symbol,
		"startTime": fmt.Sprintf("%d", startTime),
		"limit":     "100",
	}
	err := utils.MakeGetRequest(url, nil, params, 5*time.Minute, &resp)
	if err != nil {
		return []TickerFundingRate{}, err
	}
	var results []TickerFundingRate
	for _, r := range resp {
		results = append(results, TickerFundingRate{
			FundingTime: r.FundingTime,
			Symbol:      r.Symbol,
			FundingRate: r.FundingRate,
		})
	}
	return results, nil
}
