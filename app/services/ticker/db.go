package ticker

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	log "github.com/sirupsen/logrus"
)

type DB struct {
	Gorm *gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&Dataset{},
	)

	d := &DB{Gorm: db}

	return d, nil
}

func FatalWorker(timeout time.Duration, d DB) {
	log.Info("Run PG RestoreConnectionWorker")
	for {
		db, err := d.Gorm.DB()
		if err != nil {
			panic("PG is not available now")
		}

		dbWriteErr := db.Ping()
		if dbWriteErr != nil {
			panic("PG is not available now")
		}
		time.Sleep(timeout)
	}
}

func (db DB) AddDataset(dataset []Dataset) error {
	for _, b := range dataset {
		err := db.Gorm.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(&b, 200).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (db DB) GetLastTimeDataset(symbol, name string) int64 {
	var lastTime int64
	err := db.Gorm.Model(&Dataset{}).Select("unix_time").Where("symbol= AND name=", symbol, name).Order("unix_time DESC").First(&lastTime).Error
	if err != nil {
		return StartTimeMillisec
	}
	return lastTime
}
