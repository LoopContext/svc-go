package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/loopcontext/svc-go/database/models"
	"gorm.io/gorm"
)

func initialMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "INITIAL_MIGRATION",
		Migrate: func(tx *gorm.DB) (err error) {
			return tx.AutoMigrate(
				models.User{},
				models.UserAPIKey{},
				models.Client{},
			)
		},
		Rollback: func(tx *gorm.DB) (err error) {
			if err := tx.Migrator().DropTable(tx.Model(&models.User{}).Name()); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(tx.Model(&models.UserAPIKey{}).Name()); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(tx.Model(&models.Client{}).Name()); err != nil {
				return err
			}

			return
		},
	}
}
