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
	"github.com/ggrrrr/bui_lib/token/sign"
	"github.com/golang-jwt/jwt/v4"
)

func LoginUserRequest(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	body, _ := ioutil.ReadAll(r.Body)
	var auth AuthReq
	err := json.Unmarshal(body, &auth)
	if err != nil {
		log.Printf("errer %v\n", err)
		api.ResponseError(w, 400, "Asd", err)
		return
	}
	if auth.Email == "" {
		api.ResponseError(w, 400, "bad request", fmt.Errorf(""))
		return
	}
	if auth.Password == "" {
		api.ResponseError(w, 400, "bad request", fmt.Errorf(""))
		return
	}

	ua := api.GetUserAgent(r.Context())
	userPasswd, ok := passwd.VerifyUserPasswd(r.Context(), auth.Email, auth.Password)
	log.Printf("passwd.VerifyUserPasswd: namespace: %+v %+v %+v", ua, userPasswd, ok)
	// time.Sleep(2 * time.Second)
	parseLoginReulst(w, r.Context(), ok, userPasswd)

}

func parseLoginReulst(w http.ResponseWriter, ctx context.Context, ok *controlers.AuthError, userPasswd *models.UserPasswd) {
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
		api.ResponseErrorUnauthorized(w, fmt.Errorf("wrong email"))
		return
	}
	if ok.Result == controlers.AUTH_NOK {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("wrong email/pass"))
		return
	}
	if ok.Result == controlers.AUTH_OK {
		sendNewToken(w, ctx, userPasswd)
		return
	}
	api.ResponseError(w, 500, "unable to login", fmt.Errorf("unkown error"))
}

func sendNewToken(w http.ResponseWriter, ctx context.Context, userPasswd *models.UserPasswd) {

	tokenClaims := token.ApiClaims{
		Roles: "system",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userPasswd.Email,
		},
	}
	t, err := sign.SignKey(tokenClaims, ctx)
	if err != nil {
		log.Printf("Unable to sign token %v", err)
		api.ResponseError(w, 500, "Unable to sign token", err)
		return
	}
	out := AuthRes{
		Email: userPasswd.Email,
		Token: t,
		Attr:  userPasswd.Attr,
	}

	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		log.Printf("errer %v\n", err)
	}
}
