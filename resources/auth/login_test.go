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
	passwdC "github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_api_login/models/passwd"
	"github.com/ggrrrr/bui_api_login/resources/auth"
	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/ggrrrr/bui_lib/token/sign"
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

	newEmail1 := "asd@asd.com"
	newPass := "asdasd"
	passwdC.CreateUserPasswd(root, &passwd.UserPasswd{Email: newEmail1, Passwd: newPass, Enabled: true})

	loginTest1 := `{"email":"ggrrrr@gmail.com"}`
	res1, data1, err := runCall(loginTest1)
	if err != nil {
		t.Fatal(err)
	}
	if res1.StatusCode != 400 {
		t.Fatalf("must be error")
	}
	t.Logf("body:%s: err: %v", data1, err)

	loginTest2 := fmt.Sprintf(`{"email":"%v","password":"%v"}`, newEmail1, newPass)
	res2, data2, err := runCall(loginTest2)
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != 200 {
		t.Fatalf("must be error")
	}
	t.Logf("body:%s: err: %v", data2, err)

}
