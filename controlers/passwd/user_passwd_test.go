package passwd_test

import (
	"context"
	"testing"

	controlers "github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_api_login/models"
	"github.com/ggrrrr/bui_lib/db"
)

func TestPasswd1(t *testing.T) {
	var err error
	email := "asd@asd.com"
	plainPass := "asd"
	namespaces := []string{"bui", "localhost"}

	ctx := context.Background()
	// t.Setenv(db.DB_CLUSTER, "127.0.0.1")
	// t.Setenv(db.DB_KEYSPACE, "test")
	db.Configure()
	db.Connect()
	db.CreateSchema("passwd")
	defer db.Shutdown()
	pass := models.UserPasswd{Email: email, Passwd: plainPass, Enabled: true, Namespaces: namespaces}
	err = controlers.CreateUserPasswd(ctx, &pass)
	if err != nil {
		t.Error(err)
	}

	ok, _, err := controlers.VerifyUserPasswd(ctx, email, "plainPass")
	if ok != controlers.AUTH_NOK {
		t.Errorf("pass dont match")
	}

	ok, _, err = controlers.VerifyUserPasswd(ctx, "email", plainPass)
	if ok != controlers.AUTH_NOT_FOUND {
		t.Errorf("pass dont match")
	}

	ok, _, err = controlers.VerifyUserPasswd(ctx, email, plainPass)
	if ok != controlers.AUTH_OK {
		t.Errorf("pass dont match")
	}
}
