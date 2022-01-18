package passwd_test

import (
	"context"
	"testing"

	"github.com/ggrrrr/bui_api_login/controlers"
	"github.com/ggrrrr/bui_api_login/controlers/passwd"
	"github.com/ggrrrr/bui_api_login/models"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
)

func createUserPass(t *testing.T, ctx context.Context, pass *models.UserPasswd) {
	ok := passwd.CreateUserPasswd(ctx, pass)
	if ok.Err != nil {
		t.Error(ok.Err)
	}

}

func clean(t *testing.T) {
	err := db.Session.ExecStmt("DELETE from test.user_passwd where email='asd@asd.com'")
	if err != nil {
		t.Error(err)
	}
}

func TestPasswd1(t *testing.T) {
	email := "asd@asd.com"
	plainPass := "asd"
	newPlainPass := "123123"
	namespaces := []string{"bui", "localhost"}

	ctx := context.Background()
	// t.Setenv(db.DB_CLUSTER, "127.0.0.1")
	// t.Setenv(db.DB_KEYSPACE, "test")
	db.Configure()
	db.Connect()
	db.CreateSchema("passwd")
	defer db.Shutdown()
	clean(t)
	pass := models.UserPasswd{Email: email, Passwd: plainPass, Enabled: false, Namespaces: namespaces}
	createUserPass(t, ctx, &pass)

	_, ok := passwd.VerifyUserPasswd(ctx, email, "plainPass")
	if ok.Result != controlers.AUTH_LOCKED {
		t.Errorf("pass must locked")
	}
	clean(t)
	pass = models.UserPasswd{Email: email, Passwd: plainPass, Enabled: true, Namespaces: namespaces}
	createUserPass(t, ctx, &pass)

	_, ok = passwd.VerifyUserPasswd(ctx, "email", plainPass)
	if ok.Result != controlers.AUTH_NOT_FOUND {
		t.Errorf("pass dont match")
	}

	_, ok = passwd.VerifyUserPasswd(ctx, email, plainPass)
	if ok.Result != controlers.AUTH_OK {
		t.Errorf("pass dont match")
	}
	ok = passwd.ChangeUserPasswd(ctx, email, plainPass, newPlainPass)
	if ok.Result != controlers.AUTH_OK {
		t.Errorf("change pass did not worked")
	}
	_, ok = passwd.VerifyUserPasswd(ctx, email, newPlainPass)
	if ok.Result != controlers.AUTH_OK {
		t.Errorf("pass dont match")
	}
	_, ok = passwd.VerifyUserPasswd(ctx, email, plainPass)
	if ok.Result != controlers.AUTH_NOK {
		t.Errorf("pass dont match %v", ok)
	}
	clean(t)

}
