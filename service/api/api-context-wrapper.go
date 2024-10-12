package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler, authorize bool) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		if authorize {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				ctx.Logger.Error("Authorization-header: missing")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authFields := strings.Fields(authHeader)
			if len(authFields) != 2 || authFields[0] != "Bearer" {
				ctx.Logger.Error("Authorization-header: bad formatting")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userId, err := strconv.ParseUint(authFields[1], 10, 64)
			if err != nil {
				ctx.Logger.WithError(err).Error("Authorization-header: userId not valid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			isuid, err := rt.db.CheckUserId(userId)
			if err != nil {
				ctx.Logger.WithError(err).Error("Authorization-header: error checking UserId existence")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !isuid {
				ctx.Logger.WithField("fakeUID", isuid).Error("Authorization-header: invalid User-ID")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Authorised
			ctx.Logger.Debug("Authorization-header: valid")
			ctx.Uid = userId
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
