package helpers

import (
	"net/http"

	"github.com/loopcontext/svc-go/services"

	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/utils"
	"github.com/loopcontext/svc-go/utils/logger"
	"github.com/loopcontext/svc-go/web/models/responses"
)

type ResponseHelper interface {
	Error(w http.ResponseWriter, r *http.Request, statusCode int, err error)
	Send(w http.ResponseWriter, r *http.Request, statusCode int, obj interface{})
}

type ResponseHelperSvc struct {
	log         logger.Logger
	catalog     msgcat.MessageCatalog
	contentType string
}

// NewResponseHelperSvc set the contentType for the
func NewResponseHelperSvc(log logger.Logger, catalog msgcat.MessageCatalog, contentType string) ResponseHelper {
	return &ResponseHelperSvc{
		log:         log,
		catalog:     catalog,
		contentType: contentType,
	}
}

func (svc *ResponseHelperSvc) Send(w http.ResponseWriter, r *http.Request, statusCode int, obj interface{}) {
	data, err := utils.ToJSON(svc.buildResponse(r, obj))
	if err != nil {
		svc.Error(w, r, http.StatusInternalServerError, err)
	}
	svc.writeResponse(w, statusCode, data)
}

func (svc *ResponseHelperSvc) Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	var data []byte
	var errdata error
	switch svc.contentType {
	case "application/json":
		fallthrough
	default:
		data, errdata = utils.ToJSON(svc.buildResponseError(r, err))
	}
	if errdata != nil {
		svc.log.SendError(errdata)
	}
	svc.writeResponse(w, statusCode, data)
}

func (svc *ResponseHelperSvc) writeResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", svc.contentType)
	if _, err := w.Write(data); err != nil {
		svc.log.SendError(err)
	}
}

func (svc *ResponseHelperSvc) buildResponse(r *http.Request, resp interface{}) (mr *responses.Response) {
	msg := svc.catalog.GetMessageWithCtx(r.Context(), services.MsgCodeOK)
	return &responses.Response{
		ResponseBase: responses.ResponseBase{
			Code:    services.MsgCodeOK,
			Message: msg.ShortText,
			Details: msg.LongText,
		},
		Response: resp,
	}
}

func (svc *ResponseHelperSvc) buildResponseError(r *http.Request, err error) (mre *responses.ResponseError) {
	if _, ok := err.(*msgcat.DefaultError); !ok {
		err = svc.catalog.GetErrorWithCtx(r.Context(), 1, err.Error())
	}
	cde := err.(*msgcat.DefaultError)

	return &responses.ResponseError{
		ResponseBase: responses.ResponseBase{
			Code:    cde.ErrorCode(),
			Message: cde.GetShortMessage(),
			Details: cde.GetLongMessage(),
		},
	}
}
