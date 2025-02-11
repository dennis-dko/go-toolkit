package errorhandler

import (
	"errors"
	"net/http"
)

var (
	ErrAuthFailed             = errors.New("access denied to this resource")
	ErrPermFailed             = errors.New("invalid permissions to this resource")
	ErrBindingFailed          = errors.New("cannot bind the request data")
	ErrValidationFailed       = errors.New("cannot validate the request data")
	ErrDocumentNotFound       = errors.New("cannot find the document")
	ErrDocumentsNotFound      = errors.New("cannot find all documents")
	ErrDocumentNotCreate      = errors.New("cannot create the document")
	ErrDocumentNotUpdate      = errors.New("cannot update the document")
	ErrDocumentNotDelete      = errors.New("cannot delete the document")
	ErrMultipleDocumentsFound = errors.New("find multiple documents, but only one was expected")
	ErrRequestFailed          = errors.New("request failed")
	ErrRequestsLimitExceeded  = errors.New("limit of requests exceeded")
	ErrInactivityTimeout      = errors.New("inactivity timeout reached")
)

func NewErrorStatusCodeMaps() map[error]int {
	var errorStatusCodeMaps = make(map[error]int)
	errorStatusCodeMaps[ErrPermFailed] = http.StatusForbidden
	errorStatusCodeMaps[ErrAuthFailed] = http.StatusUnauthorized
	errorStatusCodeMaps[ErrBindingFailed] = http.StatusBadRequest
	errorStatusCodeMaps[ErrValidationFailed] = http.StatusBadRequest
	errorStatusCodeMaps[ErrDocumentNotFound] = http.StatusNotFound
	errorStatusCodeMaps[ErrDocumentsNotFound] = http.StatusNotFound
	errorStatusCodeMaps[ErrMultipleDocumentsFound] = http.StatusConflict
	errorStatusCodeMaps[ErrRequestsLimitExceeded] = http.StatusTooManyRequests
	return errorStatusCodeMaps
}
