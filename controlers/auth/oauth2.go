package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_api_login/controlers"
	models "github.com/ggrrrr/bui_api_login/models/passwd"
	"github.com/ggrrrr/bui_api_login/models/requests"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
	"golang.org/x/oauth2"
)

type AuthVerify struct {
	State       string
	Code        string
	Provider    string
	RedirectURL string `json:"redirect_url"`
}

type AuthProfile struct {
	ID      string
	Email   string
	Picture string
	Attr    map[string]string
}

type ProviderConfig struct {
	ClientID        string
	ClientSecret    string
	Scopes          []string
	RedirectURL     string
	AuthURL         string
	TokenURL        string
	AuthStyle       string //: oauth2.AuthStyleInParams,
	FetchProfileURL string
}

func fetchEmailGmail(client *http.Client, url string) (*AuthProfile, error) {
	resProfile, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resProfile.Body)
	log.Printf("fetchEmailGmail.body: %v", string(body))
	type profileT struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Name          string `json:"name"`
		FirstName     string `json:"given_name"`
		LastName      string `json:"family_name"`
	}
	var pp profileT
	json.Unmarshal(body, &pp)
	log.Printf("%v: %v", "", string(body))
	// log.Printf("prfile %+v\n\n\n ", pp)
	if pp.Email == "" {
		log.Printf("%v: %v", "", string(body))
		log.Printf("prfile %+v", pp)
		return nil, fmt.Errorf("unable to find email in profile")
	}
	attr := map[string]string{
		"provider":   "google",
		"id":         pp.Id,
		"first_name": pp.FirstName,
		"last_name":  pp.LastName,
		"name":       pp.Name,
		"picture":    pp.Picture,

		// "verified_email"

	}
	return &AuthProfile{Email: pp.Email, ID: pp.Id, Picture: pp.Picture, Attr: attr}, nil
}

func fetchEmailFacebook(client *http.Client, url string) (*AuthProfile, error) {
	resProfile, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resProfile.Body)
	log.Printf("fetchEmailFacebook.body: %v", string(body))
	type profileT struct {
		Id      string `json:"id"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				Url string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	var pp profileT
	json.Unmarshal(body, &pp)
	if pp.Email == "" {
		log.Printf("%v: %v", "", string(body))
		log.Printf("prfile %+v", pp)
		return nil, fmt.Errorf("unable to find email in profile")
	}
	attr := map[string]string{
		"provider": "facebook",
		"id":       pp.Id,
		"picture":  pp.Picture.Data.Url,
		// "first_name": pp.FirstName,
		// "last_name":  pp.LastName,
		// "name":       pp.Name,
		// "verified_email"

	}

	return &AuthProfile{Email: pp.Email, ID: pp.Id, Picture: pp.Picture.Data.Url, Attr: attr}, nil
}

func codeExchange(ctx context.Context, conf oauth2.Config, code string) (*http.Client, error) {
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("codeExchange: %v", err.Error())
		return nil, err
	}
	return conf.Client(ctx, tok), nil

}

func VerifyOauthCode(ctx context.Context, auth *AuthVerify) (*models.UserPasswd, *controlers.AuthError) {
	log.Printf("VerifyOauthCode")
	_, ok := fetchProfileFunc[auth.Provider]
	if !ok {
		log.Printf("auth config %+v", providers)
		return nil, controlers.ErrorStringf("auth func provider(%v) not defined", auth.Provider)
	}
	log.Printf("VerifyOauthCode...")
	p, ok := providers[auth.Provider]
	log.Printf("VerifyOauthCode: %v..", p)
	if !ok {
		log.Printf("auth config %+v", providers)
		return nil, controlers.ErrorStringf("auth provider(%v) not defined", auth.Provider)
	}
	conf := oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Scopes:       p.Scopes,
		// RedirectURL:  "http://localhost:8080./callback",
		RedirectURL: auth.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL: p.AuthURL,
			// TokenURL: "https://oauth2.googleapis.com/token",
			// TokenURL: "https://www.googleapis.com/oauth2/v4/token:",
			TokenURL:  p.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	authClient, err := codeExchange(ctx, conf, auth.Code)
	if err != nil {
		return nil, controlers.ErrorStringf("unable to fetch token %v", err)
	}
	profile, err := fetchProfileFunc[auth.Provider](authClient, p.FetchProfileURL)
	if err != nil {
		return nil, controlers.ErrorStringf("unable to fetch profile %v", err)
	}
	log.Printf("VerifyOauthCode: provider: %v, state: %v, email: %v", auth.Provider, auth.State, profile.Email)
	userPasswd, err := models.Get(db.Session, profile.Email)
	if err != nil {
		return nil, controlers.Error(err)
	}
	if userPasswd == nil {
		// profile.
		newR := requests.Request{Email: profile.Email, Enabled: false, Attr: profile.Attr}
		newR.Insert(db.Session)
		return nil, controlers.New(controlers.AUTH_NOT_FOUND, nil)
	}
	if !userPasswd.Enabled {
		return nil, controlers.New(controlers.AUTH_LOCKED, nil)
	}
	if profile.Picture != "" {
		log.Printf("asdasd %+v", userPasswd)
		if userPasswd.Attr["picture"] == "" {
			userPasswd.Attr["picture"] = profile.Picture
			log.Printf("update passwdAttr for: :%v %+v", userPasswd.Email, userPasswd.UpdateAttr(db.Session))
		}

		userPasswd.Attr["picture"] = profile.Picture
	}
	log.Printf("asdasd %+v", userPasswd)
	return userPasswd, controlers.NewOK()
	// return AUTH_ERR, nil, fmt.Errorf("shit")
}
