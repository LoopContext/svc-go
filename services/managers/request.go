package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/loopcontext/svc-go/services"

	"github.com/golang/gddo/httputil/header"
	"github.com/loopcontext/msgcat"
	"github.com/loopcontext/svc-go/utils/logger"
)

type RequestHelper interface {
	DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) (error, int)
}

type RequestHelperSvc struct {
	log     logger.Logger
	catalog msgcat.MessageCatalog
}

func NewRequestHelperSvc(log logger.Logger, catalog msgcat.MessageCatalog) RequestHelper {
	return &RequestHelperSvc{
		log:     log,
		catalog: catalog,
	}
}

func (svc *RequestHelperSvc) DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) (error, int) {
	ctx := r.Context()
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperNotJSON), http.StatusUnsupportedMediaType
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperBadlyFormedAtPos, syntaxError.Offset),
				http.StatusBadRequest

		case errors.Is(err, io.ErrUnexpectedEOF):
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperBadlyFormed), http.StatusBadRequest

		case errors.As(err, &unmarshalTypeError):
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperInvalidValue,
				unmarshalTypeError.Field, unmarshalTypeError.Offset), http.StatusBadRequest

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperUnknownField, fieldName), http.StatusBadRequest

		case errors.Is(err, io.EOF):
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperReqBodyEmpty), http.StatusBadRequest

		case err.Error() == "http: request body too large":
			return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperLimitSize), http.StatusRequestEntityTooLarge

		default:
			return err, http.StatusInternalServerError
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeReqHelperLimit1Obj), http.StatusBadRequest
	}

	return nil, http.StatusOK
}
