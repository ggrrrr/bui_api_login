package auth_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// c "github.com/ggrrrr/bui_api_login/contolers/auth"
	c "github.com/ggrrrr/bui_api_login/controlers/auth"
	"github.com/ggrrrr/bui_api_login/resources/auth"
	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/ggrrrr/bui_lib/token/sign"
	"github.com/golang-jwt/jwt/v4"
)

func runCall1(jwt, body string) (*http.Response, string, error) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))

	req.Header.Add(string(api.HTTP_CT_AUTH), fmt.Sprintf("%s %s", api.HTTP_AUTH_BEARER, jwt))
	w := httptest.NewRecorder()
	auth.ActivateRequest(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	return res, string(data), err

}

func runCall2(jwt, body string) (*http.Response, string, error) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))

	req.Header.Add(string(api.HTTP_CT_AUTH), fmt.Sprintf("%s %s", api.HTTP_AUTH_BEARER, jwt))
	w := httptest.NewRecorder()
	auth.ListRequest(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	return res, string(data), err

}

func TestReq(t *testing.T) {
	root1 := context.Background()
	root := api.SetUserAgent(root1, api.UserAgent{})
	var err error
	err = cassandra.Configure()
	if err != nil {
		log.Fatal(err)
	}
	session, err := cassandra.Connect()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer session.Close()

	if token.Configure() != nil {
		t.Fatal("token config")
	}

	if sign.Configure() != nil {
		t.Fatalf(err.Error())
	}

	c.Configure()

	tokenClaims := token.ApiClaims{
		Roles: "system",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "asd",
		},
	}
	jwt, err := sign.SignKey(tokenClaims, root)
	if err != nil {
		t.Fatalf("jwt %v", err)
	}
	runCall1(jwt, `{"email":"asd@asd.com"}`)
	a1, a2, err := runCall2(jwt, "")
	t.Logf("a1: %v a2: %v err %v", a1, a2, err)
}
