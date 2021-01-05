package repo

import (
	"context"
	"errors"

	"github.com/loopcontext/svc-go/database"
	"github.com/loopcontext/svc-go/database/models"
	"github.com/loopcontext/svc-go/services"
	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/utils/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	TblUsers       = "users"
	TblUserAPIKeys = "user_api_keys"
)

type AuthRepo interface {
	CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error)
}

type AuthRepoSvc struct {
	db      *database.DB
	log     logger.Logger
	catalog msgcat.MessageCatalog
}

func NewAuthRepoSvc(db *database.DB, log logger.Logger, catalog msgcat.MessageCatalog) AuthRepo {
	return &AuthRepoSvc{
		db:      db,
		log:     log,
		catalog: catalog,
	}
}

func (svc *AuthRepoSvc) CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error) {
	uak := &models.UserAPIKey{APIKey: apiKey}
	err = svc.db.Model(uak).Preload(clause.Associations).
		Preload("User").Preload("User.Client").Where("api_key = ?", apiKey).First(&uak).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, gorm.ErrRecordNotFound, services.MsgCodeUnauthorized, TblUserAPIKeys)
	} else if err != nil {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, err, services.MsgCodeDBUnexpectedErr, err.Error())
	}
	user = &uak.User

	return
}
