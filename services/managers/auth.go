package helpers

import (
	"context"

	"github.com/loopcontext/svc-go/database/models"
	"github.com/loopcontext/svc-go/services"
	"github.com/loopcontext/svc-go/utils"
	"github.com/loopcontext/msgcat"
	"github.com/valyala/fasthttp"
)

type AuthHelper interface {
	GetCurrentUser(ctx context.Context) (currentUser *models.User, err error)
}

type AuthHelperSvc struct {
	catalog msgcat.MessageCatalog
}

func NewAuthHelperSvc(catalog msgcat.MessageCatalog) AuthHelper {
	return &AuthHelperSvc{
		catalog: catalog,
	}
}

func (svc *AuthHelperSvc) GetCurrentUser(ctx context.Context) (currentUser *models.User, err error) {
	if ctx, ok := ctx.(*fasthttp.RequestCtx); ok {
		currentUser = ctx.UserValue(string(utils.CurrrentUserCtxKey)).(*models.User)
	} else {
		currentUser = ctx.Value(utils.CurrrentUserCtxKey).(*models.User)
	}
	if currentUser == nil {
		err = svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeHelperCurrentUserNotFound)
	}

	return
}
