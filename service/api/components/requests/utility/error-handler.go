package utility

import (
	"net/http"

	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func LogError(errorText string, statusCode int, responseWriter http.ResponseWriter, ctx reqcontext.RequestContext) {
	ctx.Logger.Error(errorText)
	responseWriter.WriteHeader(statusCode)
}

func LogWithError(errorText string, statusCode int, requestErrors error, responseWriter http.ResponseWriter, ctx reqcontext.RequestContext) {
	ctx.Logger.WithError(requestErrors).Error(errorText)
	responseWriter.WriteHeader(statusCode)
}

func LogWithField(errorText string, statusCode int, fieldName string, value interface{}, responseWriter http.ResponseWriter, ctx reqcontext.RequestContext) {
	ctx.Logger.WithField(fieldName, value).Error(errorText)
	responseWriter.WriteHeader(statusCode)
}
