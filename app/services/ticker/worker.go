package ticker

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const StartTimeMillisec = 1653004800000 // Friday, 20 May 2022 00:00:00

type Worker struct {
	client *BinanceClient
	db     *DB
}

func Init(client *BinanceClient, db *DB) *Worker {
	go FatalWorker(10*time.Second, *db)
	return &Worker{
		client: client,
		db:     db,
	}
}

func (w Worker) AddOperation(c *cron.Cron, updateTime string, f func()) {
	spec := fmt.Sprintf("@every %s", updateTime)

	if _, err := c.AddFunc(spec, f); err != nil {
		log.Fatal(err)
	}
}

func (w Worker) FetchAndSavePriceVolDataset() {
	symbol := "KNCUSDT"
	lastTime := w.db.GetLastTimeDataset(symbol, symbol+"-vol")
	if time.Now().UnixMilli()-lastTime-int64(5*time.Millisecond) < 0 {
		return
	}
	log.Printf("lastTime of dataset in database = %d", lastTime)
	tickerPrice, err := w.client.GetTickerPrice(symbol, "5m", lastTime+1, lastTime+300000*500)
	if err != nil {
		return
	}
	priceDataset := BuildPriceDataset(tickerPrice)
	if err := w.db.AddDataset(priceDataset); err != nil {
		return
	}
	volDataset := BuildVolDataSet(tickerPrice)
	if err := w.db.AddDataset(volDataset); err != nil {
		return
	}
}

func (w Worker) FetchAndSaveFundingRate() {
	symbol := "KNCUSDT"
	lastTime := w.db.GetLastTimeDataset(symbol, symbol+"-fundingRate")
	if time.Now().UnixMilli()-lastTime-int64(5*time.Millisecond) < 0 {
		return
	}
	log.Printf("lastTime of fundingRate in database = %d", lastTime)
	fundingRates, err := w.client.GetListFundingRate(symbol, lastTime+1)
	if err != nil {
		return
	}
	priceDataset := BuildFundingRateModel(fundingRates)
	if err := w.db.AddDataset(priceDataset); err != nil {
		return
	}
}
