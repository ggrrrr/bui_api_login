package passwd

import (
	"fmt"
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

/*
 */
type UserPasswd struct {
	Email       string `json:"email"`
	Passwd      string
	Enabled     bool              `json:"enabled"`
	Attr        map[string]string `json:"attr"`
	Namespaces  []string          `json:"namespaces"`
	CreatedTime time.Time         `json:"created_time"`
}

var metadata = table.Metadata{
	Name:    "user_passwd",
	Columns: []string{"email", "enabled", "created_time", "passwd", "namespaces", "attr"},
	PartKey: []string{"email"},
	// SortKey: []string{"created_time"},
}

var userPasswdTable = table.New(metadata)

func (o *UserPasswd) Insert(session *gocqlx.Session) error {
	o.CreatedTime = time.Now()
	q := session.Query(userPasswdTable.Insert()).BindStruct(o)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("unable to insert loginPasswd(%v):%+v", o.Email, err)
	}
	return nil
}

func Get(session *gocqlx.Session, email string) (*UserPasswd, error) {
	var userPasswd []UserPasswd
	q := session.Query(userPasswdTable.Select()).BindMap(qb.M{"email": email})
	if err := q.SelectRelease(&userPasswd); err != nil {
		return nil, fmt.Errorf("unable to get loginPasswd(%v):%+v", email, err)

	}
	if len(userPasswd) == 0 {
		return nil, nil
	}
	out := userPasswd[0]
	if out.Attr == nil {
		out.Attr = map[string]string{}
	}
	return &out, nil
}
func List(session *gocqlx.Session) ([]UserPasswd, error) {
	qb1 := qb.Select(metadata.Name)
	qb1 = qb1.Columns(metadata.Columns...)
	q := qb1.Query(*session)
	if q.Err() != nil {
		return nil, q.Err()
	}
	iter := q.Iter()
	if iter == nil {
		log.Printf("List: %v", q.Err())
		return nil, q.Err()

	}
	var out []UserPasswd
	var req UserPasswd
	for iter.StructScan(&req) {
		req.Passwd = ""
		out = append(out, req)
	}
	return out, nil
}

func (o *UserPasswd) UpdatePasswrd(session *gocqlx.Session) error {
	var userPasswd []UserPasswd
	q := session.Query(userPasswdTable.Update("passwd")).BindMap(qb.M{"email": o.Email, "passwd": o.Passwd})
	if err := q.SelectRelease(&userPasswd); err != nil {
		return fmt.Errorf("unable to get loginPasswd(%v):%+v", o.Email, err)

	}
	return nil
}

func (o *UserPasswd) UpdateAttr(session *gocqlx.Session) error {
	var userPasswd []UserPasswd
	q := session.Query(userPasswdTable.Update("attr")).BindMap(qb.M{"email": o.Email, "attr": o.Attr})
	if err := q.SelectRelease(&userPasswd); err != nil {
		return fmt.Errorf("unable to get loginPasswd(%v):%+v", o.Email, err)

	}
	return nil
}
