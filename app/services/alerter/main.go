package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
		timeframe := time.Minute * 10
		startTime := time.Now().Add(-timeframe)
		data, _ := getDataset(db, set, startTime)
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

func getDataset(db *gorm.DB, set string, startTime time.Time) ([]Dataset, error) {
	data := []Dataset{}
	err := db.Model(&Dataset{}).Where("set = ? AND timestamp >= ?", "knc", startTime).Find(&data).Error
	return data, err
}

type Abnormal struct {
	name        string
	level       int
	extremeData Dataset
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

	abnormals := []Abnormal{}
	for k, v := range collection {
		logger := log.With("dataset", set, "series", k)
		// log.Infow("series", "name", k, "data point", len(v))

		values := []float64{}
		for _, d := range v {
			values = append(values, d.Value)
		}

		// level 2 mean 95% | level 3 mean 99.76%
		level := 3
		o, _ := DetectOutliers(values, level)
		if len(o.Extreme) > 0 {
			median, _ := stats.Median(values)
			logger.Infow("ABNORMAL DETECTED", "extreme values", o.Extreme, "median", median)

			// find the first extreme dataset
			var extremeData Dataset
			for _, d := range v {
				if d.Value == o.Extreme[0] {
					extremeData = d
					break
				}
			}

			abnormals = append(abnormals, Abnormal{k, level, extremeData})
		} else {
			logger.Info("looks normal")
		}

		if len(abnormals) == len(collection) {
			var notes []string
			var extremeIDs []string
			for _, a := range abnormals {
				notes = append(notes, fmt.Sprintf("%s %f %v", a.name, a.extremeData.Value, a.extremeData.Timestamp))
				extremeIDs = append(extremeIDs, fmt.Sprintf("%d", a.extremeData.ID))
			}
			alert := Alert{
				Set:        set,
				Level:      fmt.Sprintf("S.O.S-level %d- %d/%d metrics", level, len(abnormals), len(collection)),
				Note:       strings.Join(notes, "|"),
				ExtremeIDs: strings.Join(extremeIDs, "-"),
				// Note:  fmt.Sprintf("extreme: [%v] | median: [%f]", o.Extreme, median),
			}
			createAlert(log, db, alert)
		}

		// o, _ := stats.QuartileOutliers(values)
		// if len(o.Extreme) > 0 {
		// 	logger.Infow("ABNORMAL DETECTED", "extreme values", o.Extreme, "mild", o.Mild)

		// 	median, _ := stats.Median(values)
		// 	alert := Alert{
		// 		Set:   set,
		// 		Name:  k,
		// 		Level: "S.O.S-0",
		// 		Note:  fmt.Sprintf("extreme: [%v] | median: [%f]", o.Extreme, median),
		// 	}
		// 	createAlert(log, db, alert)
		// } else {
		// 	logger.Info("looks normal")
		// }

	}
}

type Outliers struct {
	Level   float64
	Extreme []float64
}

// If the distribution of nums items is assumed to be normal one,
// we can treat a value being anomaly if it's beyond
// [mean - k * sigma..mean + k * sigma] range
// (sigma stands for the standard deviation),
// where k is typically 2 (95%), 3 (99.76%), sometimes even 5.
func DetectOutliers(input []float64, k int) (Outliers, error) {
	var level float64
	var extreme []float64

	mean, err := stats.Mean(input)
	if err != nil {
		return Outliers{}, err
	}

	sigma, err := stats.StandardDeviation(input)
	if err != nil {
		return Outliers{}, err
	}

	for _, v := range input {
		if v < (mean-float64(k)*sigma) || v > (mean+float64(k)*sigma) {
			extreme = append(extreme, v)
		}
	}

	// Wrap them into our struct
	return Outliers{level, extreme}, nil
}

type Alert struct {
	gorm.Model

	Set        string
	Level      string
	Note       string
	ExtremeIDs string `gorm:"uniqueIndex"`
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Alert{})
}

func createAlert(log *zap.SugaredLogger, db *gorm.DB, alert Alert) {
	db.Create(&alert)
}
