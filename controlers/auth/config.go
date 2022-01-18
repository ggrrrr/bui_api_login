package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ggrrrr/bui_lib/config"
	"github.com/spf13/viper"
)

const (
	REDIRECT_URL  = "oauth2.redirect.url"
	CLIENT_ID     = "oauth2.%s.client.id"
	CLIENT_SECRET = "oauth2.%s.client.secret"
	SCOPES        = "oauth2.%s.scopes"
	AUTH_URL      = "oauth2.%s.auth.url"
	TOKEN_URL     = "oauth2.%s.token.url"
	PROFILE_URL   = "oauth2.%s.profile.url"
)

var (
	fetchProfileFunc = map[string]func(*http.Client, string) (*AuthProfile, error){
		"google":   fetchEmailGmail,
		"facebook": fetchEmailFacebook,
	}

	// envParamsDefaults = []config.ParamValue{
	// 	{
	// 		Name:     REDIRECT_URL,
	// 		Info:     "oauth2 redirect urtl to FE",
	// 		DefValue: "http://localhost:8080",
	// 	},
	// }

	envParamsProvider = []config.ParamValue{
		{
			Name:     CLIENT_ID,
			Info:     "oauth2 %s app client id",
			DefValue: "",
		},
		{
			Name:     CLIENT_SECRET,
			Info:     "oauth2 %s app secret",
			DefValue: "",
		},
		{
			Name:     SCOPES,
			Info:     "oauth2 %s scopes",
			DefValue: "email",
		},
		{
			Name:     AUTH_URL,
			Info:     "oauth2 %s token.url",
			DefValue: "",
		},
		{
			Name:     TOKEN_URL,
			Info:     "oauth2 %s token.url",
			DefValue: "",
		},
		{
			Name:     PROFILE_URL,
			Info:     "oauth2 %s fetch profile info url",
			DefValue: "",
		},
	}

	providers = map[string]ProviderConfig{}
)

func getViper(name, group string) string {
	return viper.GetString(fmt.Sprintf(name, group))
}

func Configure() {
	// config.Configure(envParamsDefaults)

	for k := range fetchProfileFunc {
		log.Printf("config: oauth2: %v", k)
		config.ConfigureGroup(envParamsProvider, k)
	}
	for p := range fetchProfileFunc {
		clientID := getViper(CLIENT_ID, p)
		log.Printf("config: oauth2: %+v: clientID: %s", p, clientID)
		if clientID == "" {
			log.Printf("ERROR config: oauth2: %+v: : %s", p, clientID)
			continue
		}
		providers[p] = ProviderConfig{
			ClientID:     getViper(CLIENT_ID, p),
			ClientSecret: getViper(CLIENT_SECRET, p),
			Scopes:       strings.Split(getViper(SCOPES, p), ","),
			AuthURL:      getViper(AUTH_URL, p),
			TokenURL:     getViper(TOKEN_URL, p),
			// AuthStyle:       "AuthStyleInParams",
			FetchProfileURL: getViper(PROFILE_URL, p),
		}
	}

}
