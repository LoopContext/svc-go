package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, migrations ...*gormigrate.Migration) (err error) {
	e := append([]*gormigrate.Migration{
		initialMigration(),
		firstUserMigration(),
	}, migrations...)
	m := gormigrate.New(db, gormigrate.DefaultOptions, e)

	return m.Migrate()
}
