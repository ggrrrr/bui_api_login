package auth_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// c "github.com/ggrrrr/bui_api_login/contolers/auth"
	c "github.com/ggrrrr/bui_api_login/controlers/auth"
	"github.com/ggrrrr/bui_api_login/resources/auth"
	"github.com/ggrrrr/bui_lib/db/cassandra"
)

func runCall(body string) (*http.Response, string, error) {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	auth.LoginUserRequest(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	return res, string(data), err

}

func TestLogin(t *testing.T) {
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

	c.Configure()

	loginTest := `{"email":"ggrrrr@gmail.com"}`
	res, data, err := runCall(loginTest)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("must be error")
	}
	t.Logf("body:%s: err: %v", data, err)

}
