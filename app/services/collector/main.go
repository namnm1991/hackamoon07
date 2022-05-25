package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/montanaflynn/stats"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	// ======================================================
	// Construct the application logger.
	log, err := logger.New("smart-alert-collector")
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
	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		DB struct {
			User       string `conf:"default:suser"`
			Password   string `conf:"default:spassword,mask"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:smart-alert"`
			DisableTLS bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SMART-ALERT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// ======================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
	}()

	// ======================================================
	// Generate datasets

	generateDataset(db)

	return nil
}

type Dataset struct {
	gorm.Model

	Set       string
	Name      string
	Value     float64
	Timestamp time.Time
}

func generateDataset(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&Dataset{})

	// =========================================================================
	set := "knc"
	lens := 1200
	interval := time.Second * 30
	startTime := time.Now().Add(-time.Duration(lens) * interval)

	dataset := []Dataset{}
	// price from 1.8 to 2.2
	d1 := g(set, "knc_price", lens, 2, 0.2, interval, startTime)
	// vol from 9000 to 11000
	d2 := g(set, "knc_vol", lens, 10000, 1000, interval, startTime)
	// vol from 3000 to 7000
	d3 := g(set, "knc_transfer", lens, 5000, 2000, interval, startTime)
	dataset = append(dataset, d1...)
	dataset = append(dataset, d2...)
	dataset = append(dataset, d3...)

	// Save the dataset into the database
	for _, d := range dataset {
		// Create
		db.Create(&d)
	}
}

func g(set, name string, lens int, loc float64, scale float64, interval time.Duration, startTime time.Time) []Dataset {
	dataset := []Dataset{}
	data := stats.NormBoxMullerRvs(loc, scale, lens)
	for i := 0; i < lens; i++ {
		dataset = append(dataset, Dataset{
			Set:       set,
			Name:      name,
			Value:     data[i],
			Timestamp: startTime.Add(time.Duration(i) * interval),
		})
	}

	return dataset
}
