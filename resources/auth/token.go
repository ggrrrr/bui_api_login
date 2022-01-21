package auth

import (
	"log"
	"net/http"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/token"
)

func TokenVerifyRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	var err error
	t, err := api.GetAuthorizationBearer(r)
	log.Printf("TokenVerifyRequest(%v): token: %v", api.GetUserAgent(r.Context()), t)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	_, err = token.Verify(t, r.Context())
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	api.ResponseOk(w, "ok", nil)
}
