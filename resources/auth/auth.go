package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_lib/api"
)

const (
	// View your email address
	UserinfoEmailScope = "https://www.googleapis.com/auth/userinfo.email"

	// See your personal info, including any personal info you've made
	// publicly available
	UserinfoProfileScope = "https://www.googleapis.com/auth/userinfo.profile"

	// Associate you with your personal info on Google
	OpenIDScope = "openid"
)

func LoginOauth2Request(w http.ResponseWriter, r *http.Request) {
	api.SetResponseHeader(w)
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ResponseError(w, 400, "Unable to read body", err)
		return
	}
	auth := passwd.AuthVerify{}
	json.Unmarshal(body, &auth)
	log.Printf("LoginOauth2Request provider: %v, state: %v url: %s, code: %v", auth.Provider, auth.State, auth.RedirectURL, auth.Code)
	ok, userPasswd, err := passwd.VerifyOauthCode(r.Context(), &auth)

	parseLoginData(w, r.Context(), ok, userPasswd, err)
}
