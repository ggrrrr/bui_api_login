package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_api_login/controlers"
	"github.com/ggrrrr/bui_api_login/controlers/passwd"
	models "github.com/ggrrrr/bui_api_login/models/passwd"
	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/token"
)

func ChangePasswordRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	t, err := api.GetAuthorizationBearer(r)
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	jwt, err := token.Verify(t, r.Context())
	// jwt.
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}
	if err != nil {
		api.ResponseErrorUnauthorized(w, err)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var newPassowrd ChangePasswordReq
	err = json.Unmarshal(body, &newPassowrd)
	if err != nil {
		log.Printf("errer %v\n", err)
		api.ResponseError(w, 400, "Asd", err)
		return
	}

	ok := passwd.ChangeUserPasswd(r.Context(), jwt.Subject, newPassowrd.Password, newPassowrd.NewPassword)
	parsePasswordReulst(w, r.Context(), ok, nil)
}

func parsePasswordReulst(w http.ResponseWriter, ctx context.Context, ok *controlers.AuthError, userPasswd *models.UserPasswd) {
	email := ""
	enabled := false
	if userPasswd != nil {
		email = userPasswd.Email
		enabled = userPasswd.Enabled
	}
	log.Printf("parseLoginReulst: user: %v enabled: %v , err: %v", email, enabled, ok)
	if ok.Err != nil {
		log.Printf("error %v\n", ok.Err)
		api.ResponseError(w, 500, "InternalError", ok.Err)
		return
	}
	if ok.Result == controlers.AUTH_LOCKED {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("account is locked"))
		return
	}
	if ok.Result == controlers.AUTH_NOT_FOUND {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("account"))
		return
	}
	if ok.Result == controlers.AUTH_NOK {
		api.ResponseError(w, 400, "wrong password", fmt.Errorf("wrong password"))
		return
	}
	if ok.Result == controlers.AUTH_OK {
		api.ResponseOk(w, "ok", nil)
		return
	}
	api.ResponseError(w, 500, "unable to login", fmt.Errorf("unkown error"))
}
