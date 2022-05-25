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

	// Construct the application logger.
	log, err := logger.New("smart-alert-monitor")
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

	// Migrate
	migrate(db)

	// ======================================================
	// Monitoring

	ticker := time.NewTicker(5 * time.Second).C
	for range ticker {
		log.Info("checking ...")

		set := "knc"
		timeframe := time.Minute * 5
		data, _ := getDataset(db, set, timeframe)

		checkAbnormal(log, db, data, set)
	}

	return nil
}

type Dataset struct {
	gorm.Model

	Set       string
	Name      string
	Value     float64
	Timestamp time.Time
}

func getDataset(db *gorm.DB, set string, timeframe time.Duration) ([]Dataset, error) {
	data := []Dataset{}
	startTime := time.Now().Add(-timeframe)
	err := db.Model(&Dataset{}).Where("set = ? AND timestamp >= ?", "knc", startTime).Find(&data).Error
	return data, err
}

func checkAbnormal(log *zap.SugaredLogger, db *gorm.DB, data []Dataset, set string) {
	collection := map[string][]Dataset{}

	for _, d := range data {
		series, ok := collection[d.Name]
		if !ok {
			series = []Dataset{}
		}

		series = append(series, d)
		collection[d.Name] = series
	}

	for k, v := range collection {
		logger := log.With("dataset", set, "series", k)
		// log.Infow("series", "name", k, "data point", len(v))

		values := []float64{}
		for _, d := range v {
			values = append(values, d.Value)
		}
		o, _ := stats.QuartileOutliers(values)
		if len(o.Extreme) > 0 {
			logger.Infow("ABNORMAL DETECTED", "extreme values", o.Extreme, "mild", o.Mild)

			median, _ := stats.Median(values)
			alert := Alert{
				Set:   set,
				Name:  k,
				Level: "S.O.S-0",
				Note:  fmt.Sprintf("extreme: [%v] | median: [%f]", o.Extreme, median),
			}
			createAlert(log, db, alert)
		} else {
			logger.Info("looks normal")
		}
	}
}

type Alert struct {
	gorm.Model

	Set   string
	Name  string
	Level string
	Note  string
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Alert{})
}

func createAlert(log *zap.SugaredLogger, db *gorm.DB, alert Alert) {
	db.Create(&alert)
}
