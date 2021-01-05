package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/loopcontext/svc-go/database/models"
	"gorm.io/gorm"
)

func firstUserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "FIRST_CLIENT_AND_USER",
		Migrate: func(tx *gorm.DB) (err error) {
			u := &models.Client{
				Name: "LoopContext",
				VAT: "123456789-0",
				Users: []models.User{
					{
						FullName: "Admin",
						Email:    "admin@loopcontext.com",
						UserName: "test",
						Password: "test",
						APIKeys: []models.UserAPIKey{
							{}, // Generates an API key
						},
					},
				},
			}
			if err = tx.FirstOrCreate(u).Error; err != nil {
				return err
			}

			return
		},
		Rollback: func(tx *gorm.DB) (err error) {
			return
		},
	}
}
