package repo

import (
	"context"

	"github.com/loopcontext/svc-go/database"
	"github.com/loopcontext/svc-go/database/models"
	"gorm.io/gorm/clause"
)

type ClientRepo interface {
	UpsertClient(ctx context.Context, input *models.Client) (*models.Client, error)
	ListByClientID(ctx context.Context, clientID string) (list []models.Client, err error)
	Read(ctx context.Context, clientID string) (client *models.Client, err error)
}

type ClientRepoSvc struct {
	db *database.DB
}

func NewClientRepoSvc(db *database.DB) ClientRepo {
	return &ClientRepoSvc{
		db: db,
	}
}

func (svc *ClientRepoSvc) Read(ctx context.Context, clientID string) (*models.Client, error) {
	rec := &models.Client{BaseModel: models.BaseModel{ID: clientID}}
	err := svc.db.Model(rec).Preload(clause.Associations).First(rec).Error

	return rec, err
}

func (svc *ClientRepoSvc) UpsertClient(ctx context.Context, input *models.Client) (*models.Client, error) {
	err := svc.db.Model(input).FirstOrCreate(input, input).Save(input).Error
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (svc *ClientRepoSvc) ListByClientID(ctx context.Context, clientID string) (list []models.Client, err error) {
	err = svc.db.Model(list).Where("client_id = ?", clientID).Find(&list).Error

	return
}
