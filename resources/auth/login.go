package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_api_login/models"
	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/ggrrrr/bui_lib/token/sign"
	"github.com/golang-jwt/jwt"
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
	log.Printf("%s %+v", r.Method, auth.Email)

	ok, userPasswd, err := passwd.VerifyUserPasswd(r.Context(), auth.Email, auth.Password)
	// time.Sleep(2 * time.Second)
	parseLoginData(w, r.Context(), ok, userPasswd, err)

}

func parseLoginData(w http.ResponseWriter, ctx context.Context, ok passwd.AUTH_RESULT, userPasswd *models.UserPasswd, err error) {
	email := ""
	enabled := false
	if userPasswd != nil {
		email = userPasswd.Email
		enabled = userPasswd.Enabled
	}
	log.Printf("parseLoginData: %v, user: %v enabled: %v , err: %v", passwd.AUTH_INFO(ok), email, enabled, err)
	if err != nil {
		log.Printf("error %v\n", err)
		api.ResponseError(w, 500, "Asd", err)
		return
	}
	if ok == passwd.AUTH_LOCKED {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("account is locked"))
		return
	}
	if ok == passwd.AUTH_NOT_FOUND {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("wrong email"))
		return
	}
	if ok == passwd.AUTH_NOK {
		api.ResponseErrorUnauthorized(w, fmt.Errorf("wrong email/pass"))
		return
	}
	if ok == passwd.AUTH_OK {

		sendNewToken(w, ctx, userPasswd)
		return
	}
	api.ResponseError(w, 500, "unable to login", fmt.Errorf("unkown error"))

	// asd := jwt.Aasd(login.Email)

}

func sendNewToken(w http.ResponseWriter, ctx context.Context, userPasswd *models.UserPasswd) {

	tokenClaims := token.ApiClaims{
		Groups: "asdasdasd",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Subject:   fmt.Sprintf("email:%s", userPasswd.Email),
			Issuer:    api.Name(),
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
