package passwd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ggrrrr/bui_api_login/models"
	"github.com/ggrrrr/bui_lib/db"
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
	type profileT struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
	}
	var pp profileT
	json.Unmarshal(body, &pp)
	// log.Printf("%v: %v", "", string(body))
	// log.Printf("prfile %+v\n\n\n ", pp)
	if pp.Email == "" {
		log.Printf("%v: %v", "", string(body))
		log.Printf("prfile %+v", pp)
		return nil, fmt.Errorf("unable to find email in profile")
	}
	return &AuthProfile{Email: pp.Email, ID: pp.Id, Picture: pp.Picture}, nil
}

func fetchEmailFacebook(client *http.Client, url string) (*AuthProfile, error) {
	resProfile, _ := client.Get(url)
	body, _ := ioutil.ReadAll(resProfile.Body)
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
	return &AuthProfile{Email: pp.Email, ID: pp.Id, Picture: pp.Picture.Data.Url}, nil
}

func codeExchange(ctx context.Context, conf oauth2.Config, code string) (*http.Client, error) {
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return conf.Client(ctx, tok), nil

}

func VerifyOauthCode(ctx context.Context, auth *AuthVerify) (AUTH_RESULT, *models.UserPasswd, error) {
	_, ok := fetchProfileFunc[auth.Provider]
	if !ok {
		return AUTH_ERR, nil, fmt.Errorf("auth provider(%v) not defined", auth.Provider)
	}
	// make sure we have config and func for provider X
	p, ok := providers[auth.Provider]
	if !ok {
		return AUTH_ERR, nil, fmt.Errorf("auth provider(%v) not defined", auth.Provider)
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
		return AUTH_ERR, nil, fmt.Errorf("unable to fetch token %v", err)
	}
	profile, err := fetchProfileFunc[auth.Provider](authClient, p.FetchProfileURL)
	if err != nil {
		return AUTH_ERR, nil, fmt.Errorf("unable to fetch profile %v", err)
	}
	log.Printf("VerifyOauthCode: provider: %v, state: %v, email: %v", auth.Provider, auth.State, profile.Email)
	passwd, err := models.Get(db.Session, profile.Email)
	if err != nil {
		return AUTH_ERR, nil, err
	}
	if passwd == nil {
		return AUTH_NOT_FOUND, nil, nil
	}
	if !passwd.Enabled {
		return AUTH_LOCKED, nil, nil
	}
	if profile.Picture != "" {
		log.Printf("asdasd %+v", passwd)
		if passwd.Attr["picture"] == "" {
			passwd.Attr["picture"] = profile.Picture
			log.Printf("update passwdAttr for: :%v %+v", passwd.Email, passwd.UpdateAttr(db.Session))
		}

		passwd.Attr["picture"] = profile.Picture
	}
	log.Printf("asdasd %+v", passwd)
	return AUTH_OK, passwd, nil
	// return AUTH_ERR, nil, fmt.Errorf("shit")
}
