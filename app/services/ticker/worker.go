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

func (w Worker) FetchAndSaveDataset() {
	lastTime := w.db.GetLastTimeDataset()
	log.Printf("lastTime of dataset in database = %d", lastTime)
	tickerPrice, err := w.client.GetTickerPrice("KNCUSDT", "5m", lastTime, lastTime+300000*500)
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
