package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/business/sys/database"
	"github.com/ardanlabs/service/foundation/logger"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("SMART-ALERT")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// demoSendEmail()
	// demoSendWebPush()

	fmt.Println(time.Now().Unix())

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func demoSendWebPush() {
	// ======================================================
	// send an web push noti
	rand.Seed(time.Now().Unix())
	no := rand.Intn(100)
	title := fmt.Sprintf("S.O.S title %d", no)
	content := fmt.Sprintf("S.O.S content %d", no)
	sendWebPush(title, content)
}

func demoSendEmail() {
	// ======================================================
	// send an email
	emails := []string{"nam@krystal.app"}
	subject := "Welcome to Krystal SmartAlert"
	content := fmt.Sprintf("S.O.S %d", rand.Intn(100))
	sendEmail(emails, subject, content)

	// ======================================================
	// draw a spark line
	// sparkline()

	// ======================================================
	// simple calculation
	// [mean - k * sigma..mean + k * sigma] range
	// (sigma stands for the standard deviation),
	// where k is typically 2 (95%), 3 (99.76%),

	// nums := []float64{3, 5, 9, 1, 8, 6, 58, 9, 4, 10}
	// m, _ := stats.Mean(nums)
	// sd, _ := stats.StandardDeviation(nums)

	// fmt.Printf("mean [%.3f], standard deviation: [%.3f]\n", m, sd)

	// o, _ := stats.QuartileOutliers([]float64{-1000, 1, 3, 4, 4, 6, 6, 6, 6, 7, 8, 15, 18, 100})
	// fmt.Printf("%+v\n", o)
}

func run(log *zap.SugaredLogger) error {

	cfg := struct {
		conf.Version
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	// ======================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	return nil
}
