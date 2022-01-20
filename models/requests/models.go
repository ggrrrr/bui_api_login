package requests

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

type Request struct {
	Email       string            `json:"email"`
	Enabled     bool              `json:"enabled"`
	Attr        map[string]string `json:"attr"`
	CreatedTime time.Time         `json:"created_time"`
}

var requestMetadata = table.Metadata{
	Name:    "request",
	Columns: []string{"email", "enabled", "created_time", "attr"},
	PartKey: []string{"email"},
}

// var userPasswdTable = table.New(userPasswdMetadata)

func (o *Request) Insert(session *gocqlx.Session) error {
	o.CreatedTime = time.Now()
	o.Enabled = false
	insertQ := qb.Insert(requestMetadata.Name)

	insertQ.Columns(requestMetadata.Columns...)
	stmt, names := insertQ.ToCql()
	log.Printf("stmt, names %v, %v", stmt, names)
	q := session.Query(stmt, names).BindStruct(o)
	err := q.Exec()
	log.Printf("Insert err: %v", err)
	return err
}

func List(session *gocqlx.Session) ([]Request, error) {
	qb1 := qb.Select(requestMetadata.Name)
	qb1 = qb1.Columns(requestMetadata.Columns...)
	q := qb1.Query(*session)
	if q.Err() != nil {
		return nil, q.Err()
	}
	iter := q.Iter()
	if iter == nil {
		log.Printf("List: %v", q.Err())
		return nil, q.Err()

	}
	var out []Request
	var req Request
	for iter.StructScan(&req) {
		if !req.Enabled {
			out = append(out, req)
		}
	}
	return out, nil
}

func (o *Request) UpdateEnable(session *gocqlx.Session) error {
	// var userPasswd []UserPasswd

	updateQ := qb.Update(requestMetadata.Name)
	updateQ = updateQ.Set("enabled").Where(qb.Eq("email"))
	stmt, names := updateQ.ToCql()
	q := session.Query(stmt, names).BindStruct(o)
	err := q.Exec()
	log.Printf("UpdateEnable: %v", err)

	// if err := q.SelectRelease(&userPasswd); err != nil {
	// return fmt.Errorf("unable to get loginPasswd(%v):%+v", o.Email, err)
	//
	// }
	return err
}
