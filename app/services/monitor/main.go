package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/service/foundation/logger"
	"github.com/gizak/termui/v3/widgets"
	"github.com/namnm1991/hackamoon07/feed"
	"go.uber.org/zap"

	"github.com/gizak/termui/v3"
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

	// ======================================================
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	// ======================================================
	// visulization
	title := "Token Price"
	refreshCh := make(chan bool)
	go sparkline(title, feed.FetchData(), refreshCh)

	// ======================================================
	// update data by interval
	// also visualize it
	interval := 5 * time.Second
	ticker := time.NewTicker(interval).C
	uiEvents := termui.PollEvents()
	for {
		select {
		case <-ticker:
			refreshCh <- true
			go sparkline(title, feed.FetchData(), refreshCh)
		case e := <-uiEvents:
			refreshCh <- true
			if e.Type == termui.KeyboardEvent && e.ID == "q" {
				return nil
			}
		}
	}
}

func sparkline(title string, data []float64, done <-chan bool) {
	sl0 := widgets.NewSparkline()
	sl0.Data = data
	sl0.LineColor = termui.ColorGreen

	// single
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = title
	slg0.SetRect(0, 0, 20, 10)

	// ===========================================================================
	// render
	termui.Render(slg0)

	// ===========================================================================
	// wait for quit event
	for {
		select {
		case <-done:
			return
		case e := <-termui.PollEvents():
			if e.Type == termui.KeyboardEvent && e.ID == "q" {
				return
			}
		}
	}
}
