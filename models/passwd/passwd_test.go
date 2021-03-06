package passwd_test

import (
	"testing"

	"github.com/ggrrrr/bui_api_login/models/passwd"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
)

func TestPasswd1(t *testing.T) {
	var err error
	// ctx := context.Background()
	// t.Setenv(db.DB_CLUSTER, "127.0.0.1")
	// t.Setenv(db.DB_KEYSPACE, "test")
	err = db.Configure()
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	session, err := db.Connect()
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	db.CreateSchema("passwd")

	userPasswd := passwd.UserPasswd{Email: "asd@asd.com", Enabled: true, Passwd: "asd", Namespaces: []string{"localhost"}, Attr: map[string]string{"asd": "asd"}}
	err = userPasswd.Insert(session)
	if err != nil {
		t.Errorf("cant insert : %+v", err)
	}

	list, err := passwd.List(session)
	if err != nil {
		t.Errorf("cant get : %+v", err)
	}
	if len(list) < 1 {
		t.Errorf("cant get : %+v", err)
	}
	t.Logf("list: %v", list)
}
