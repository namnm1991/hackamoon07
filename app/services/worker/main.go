package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/namnm1991/hackamoon07/app/services/ticker"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	// ======================================================
	// Construct the application logger.
	log, err := logger.New("smart-alert-datasource")
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
	db, err := ticker.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
	}()

	// ======================================================
	// Generate datasets

	bClient := ticker.NewBinanceClient()
	w := ticker.Init(bClient, db)
	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	w.AddOperation(c, "5m", w.FetchAndSaveDataset)

	go c.Start()
	go w.FetchAndSaveDataset()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown worker gracefully...")
	ctx := c.Stop()
	<-ctx.Done()

	return nil
}
