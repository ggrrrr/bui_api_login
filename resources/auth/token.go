package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/token"
)

func TokenVerifyRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	var err error
	t, err := api.GetAuthorizationBearer(r)

	log.Printf("TokenVerifyRequest(%v): token: %v", r.Context().Value(api.ContextUserAgent), t)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	jwt, err := token.Verify(t, r.Context())
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	if !jwt.Valid {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("invalid JWT"))
		return
	}
	api.ResponseOk(w)
}
