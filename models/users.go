package models

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

/*
CREATE TABLE users(
namespace text,
id uuid,
name text,
email text,
first_name text,
last_name text,
phones frozen<map<text,text>>,
attr frozen<map<text,text>>,
labels frozen<map<text,text>>,
created_time timestamp,
PRIMARY KEY (namespace, id)
)

create index users_email on users ( email)

*/

type User struct {
	Namespace   string
	Email       string
	Id          gocql.UUID
	Name        string
	FirstName   string
	LastName    string
	Phones      map[string]string
	Labels      map[string]string
	Attr        map[string]string
	CreatedTime time.Time
}

var userMetadata = table.Metadata{
	Name:    "users",
	Columns: []string{"namespace", "email", "id", "name", "first_name", "last_name", "attr", "created_time"},
	PartKey: []string{"namespace"},
	SortKey: []string{"id", "email"},
}

var usersTable = table.New(userMetadata)

func (o *User) Insert(session *gocqlx.Session) error {
	o.CreatedTime = time.Now()
	q := session.Query(usersTable.Insert()).BindStruct(o)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("unable to insert user(%v):%+v", o.Email, err)
	}
	return nil
}

func (o *User) GetByEmail(session *gocqlx.Session) ([]User, error) {
	var user []User
	q := session.Query(usersTable.Select()).BindMap(qb.M{"email": o.Email, "namespace": o.Namespace})
	if err := q.SelectRelease(&user); err != nil {
		return nil, fmt.Errorf("unable to get user(%v):%+v", o.Email, err)

	}
	return user, nil
}
