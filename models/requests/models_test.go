package requests_test

import (
	"testing"

	"github.com/ggrrrr/bui_api_login/models/requests"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
)

func TestRequests(t *testing.T) {
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

	userPasswd := requests.Request{Email: "asd@asd.com", Enabled: false, Attr: map[string]string{"asd": "asd"}}
	err = userPasswd.Insert(session)
	if err != nil {
		t.Errorf("cant insert : %+v", err)
	}
}
