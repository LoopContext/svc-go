package database

import (
	"github.com/loopcontext/svc-go/database/migrations"
	"github.com/loopcontext/svc-go/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg DBConfig, log *logger.Logger) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.SendError(err)
	}
	if cfg.Automigrate {
		if err = migrations.RunMigrations(db); err != nil {
			return nil, err
		}
	}
	if cfg.Debug {
		db = db.Debug()
	}
	return &DB{
		Logger: log,
		DB:     db,
	}, nil
}
