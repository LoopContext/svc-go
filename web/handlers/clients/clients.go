package clients

import (
	"net/http"
	"path"

	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/database/repo"
	"github.com/loopcontext/svc-go/services/managers"
	"github.com/loopcontext/svc-go/utils/logger"
)

type ClientHandler interface {
	// Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	// Update(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type ClientHandlerSvc struct {
	log            logger.Logger
	catalog        msgcat.MessageCatalog
	clientRepo     repo.ClientRepo
	authHelper     helpers.AuthHelper
	requestHelper  helpers.RequestHelper
	responseHelper helpers.ResponseHelper
}

func NewClientHandlerSvc(log logger.Logger, catalog msgcat.MessageCatalog,
	clientRepo repo.ClientRepo, authHelper helpers.AuthHelper,
	requestHelper helpers.RequestHelper, responseHelper helpers.ResponseHelper) ClientHandler {
	return &ClientHandlerSvc{
		log:            log,
		catalog:        catalog,
		clientRepo:     clientRepo,
		authHelper:     authHelper,
		responseHelper: responseHelper,
		requestHelper:  requestHelper,
	}
}

func (svc *ClientHandlerSvc) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// currentUser, err := svc.authHelper.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

	// 	return
	// }
	clientID := path.Base(r.URL.Path)
	list, err := svc.clientRepo.Read(ctx, clientID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, list)
}

func (svc *ClientHandlerSvc) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authHelper.GetCurrentUser(ctx)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	list, err := svc.clientRepo.ListByClientID(ctx, currentUser.Client.ID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, list)
}
