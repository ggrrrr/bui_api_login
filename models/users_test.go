package models_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"

	"github.com/ggrrrr/bui_api_login/models"
)

func gTestUserPasswd1(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "test"

	email := "asd@asd.com"
	namespace := "localhost"

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	o := models.User{Id: gocql.TimeUUID(), Email: email, Namespace: namespace}
	err = o.Insert(&session)
	if err != nil {
		t.Errorf("cant insert : %+v", err)
	}

	user1 := models.User{Email: email, Namespace: namespace}
	asd, err := user1.GetByEmail(&session)
	t.Logf("%v %v", asd, err)
}
