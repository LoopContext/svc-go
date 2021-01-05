package middleware

import (
	"context"
	"net/http"

	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/database/repo"
	"github.com/loopcontext/svc-go/services/managers"
	"github.com/loopcontext/svc-go/utils"
	"github.com/loopcontext/svc-go/utils/logger"
)

const HeaderAPIKey = "X-API-KEY"

type AuthMiddleware interface {
	CheckUserAPIKey(h http.Handler) http.Handler
}

type AuthMiddlewareSvc struct {
	log            logger.Logger
	catalog        msgcat.MessageCatalog
	authRepo       repo.AuthRepo
	responseHelper helpers.ResponseHelper
}

func NewAuthMiddlewareSvc(log logger.Logger, catalog msgcat.MessageCatalog,
	authRepo repo.AuthRepo, responseHelper helpers.ResponseHelper) AuthMiddleware {
	return &AuthMiddlewareSvc{
		log:            log,
		authRepo:       authRepo,
		responseHelper: responseHelper,
		catalog:        catalog,
	}
}

func (svc *AuthMiddlewareSvc) CheckUserAPIKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey := r.Header.Get(HeaderAPIKey)
		if user, err := svc.authRepo.CheckUserAPIKey(r.Context(), apikey); err != nil {
			svc.responseHelper.Error(w, r, http.StatusUnauthorized, err)

			return
		} else {
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), utils.CurrrentUserCtxKey, user)))
		}
	})
}
