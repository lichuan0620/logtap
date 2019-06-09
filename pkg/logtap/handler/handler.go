package handler

import (
	"net/http"

	"github.com/lichuan0620/logtap/pkg/fieldpath"
	"github.com/lichuan0620/logtap/pkg/httputil"
	"github.com/lichuan0620/logtap/pkg/logtap"
	model "github.com/lichuan0620/logtap/pkg/model/v1alpha1"
)

const notFoundMessage = "Cannot find the requested LogTask object."

type logTapHandler struct {
	tap logtap.LogTap
}

// NewLogTapHandler returns a http.Handler that handles one LogTap instance.
func NewLogTapHandler(tap logtap.LogTap) http.Handler {
	return &logTapHandler{
		tap: tap,
	}
}

// ServeHTTP implements the http.Handler interface.
func (h *logTapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		task, err := h.getLogTask()
		httputil.WriteGetResponse(w, task, err)
	default:
		httputil.WriteGetResponse(w, nil, httputil.NewMethodNotAllowedError())
	}
}

func (h *logTapHandler) getLogTask() (*model.LogTask, error) {
	if h == nil || h.tap == nil {
		return nil, httputil.NewNotFoundError(notFoundMessage)
	}
	task := h.tap.GetTask()
	if task == nil {
		return nil, httputil.NewNotFoundError(notFoundMessage)
	}
	if err := model.ValidateLogTask(fieldpath.NewFieldPath("logTask"), task); err != nil {
		return nil, httputil.NewValidationError(err.Error())
	}
	return task, nil
}
