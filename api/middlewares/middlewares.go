package middlewares

import (
	"github.com/HasanShahjahan/go-guest/api/auth"
	"github.com/HasanShahjahan/go-guest/api/config"
	"github.com/HasanShahjahan/go-guest/api/responses"
	logging "github.com/HasanShahjahan/go-guest/api/utils"
	"strconv"
)

import (
	"net/http"
)

const (
	logTag = "[Middlewares]"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Info(logTag, "JWT authentication token :"+strconv.FormatBool(config.Config.IsJwtEnabled))
		if config.Config.IsJwtEnabled {
			err := auth.TokenValid(r)
			if err != nil {
				responses.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
				logging.Error(logTag, "Unauthorized", err)
				return
			}
		}
		next(w, r)
	}
}
