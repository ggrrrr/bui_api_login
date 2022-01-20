package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	c "github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_api_login/models/passwd"
	"github.com/ggrrrr/bui_api_login/models/requests"
	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/ggrrrr/bui_lib/token"
)

type ActivateReq struct {
	Email string `json:"email"`
}

type NewRequests struct {
	Data requests.Request `json:"data"`
}

func ActivateRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	var err error
	t, err := api.GetAuthorizationBearer(r)
	log.Printf("ActivateRequest(%v): token: %v", api.GetUserAgent(r.Context()), t)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	apiClaims, err := token.Verify(t, r.Context())
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	var auth ActivateReq
	err = json.Unmarshal(body, &auth)
	if err != nil {
		log.Printf("errer %v\n", err)
		api.ResponseError(w, 400, "Wrong data", err)
		return
	}
	if auth.Email == "" {
		log.Printf("errer email")
		api.ResponseError(w, 400, "Wrong email", err)
		return
	}
	activateReq := requests.Request{Email: auth.Email, Enabled: true}
	activateReq.UpdateEnable(cassandra.Session)
	newPasswd := passwd.UserPasswd{Email: auth.Email, Enabled: true}
	c.CreateUserPasswd(r.Context(), &newPasswd)

	log.Printf("asdActivateRequest(%v): token: %v", api.GetUserAgent(r.Context()), apiClaims)
	api.ResponseOk(w)
}

func ListRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	var err error
	t, err := api.GetAuthorizationBearer(r)
	log.Printf("ActivateRequest(%v): token: %v", api.GetUserAgent(r.Context()), t)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	apiClaims, err := token.Verify(t, r.Context())
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	asd, err := requests.List(cassandra.Session)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	_ = json.NewEncoder(w).Encode(asd)

	log.Printf("asdActivateRequest(%v): token: %v", api.GetUserAgent(r.Context()), apiClaims)
	api.ResponseOk(w)
}
