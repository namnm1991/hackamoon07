package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/service/foundation/logger"
	"github.com/namnm1991/hackamoon07/feed"
	"go.uber.org/zap"

	"github.com/montanaflynn/stats"
)

func main() {

	// Construct the application logger.
	log, err := logger.New("SMART-ALERTER")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// ======================================================
	// init the data feed
	feed := feed.NewFeed()
	defer feed.Close()

	ticker := time.NewTicker(5 * time.Second).C
	for range ticker {
		data := feed.FetchData()
		o, _ := stats.QuartileOutliers(data)

		if len(o.Extreme) > 0 {
			log.Infow("ABNORMAL DETECTED", "extreme values", o.Extreme, "mild", o.Mild)
			// TODO:
			// sending the notification here
		} else {
			log.Info("everything is normal ")
		}

		log.Info("sleeping ...")
	}

	return nil
}
