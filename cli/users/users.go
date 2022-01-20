package users

import (
	"flag"
	"fmt"

	"github.com/ggrrrr/bui_api_login/cli"
	models "github.com/ggrrrr/bui_api_login/models/passwd"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
	"golang.org/x/crypto/bcrypt"
)

const name = "users"

type some struct{}

func New() cli.CliCommand {
	return &some{}
}

func (c *some) Help() {
	fmt.Printf("%s Help", name)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func updatePasswd(email string, newPass string) {
	hashPass, _ := hashPassword(newPass)
	passwd, err := models.Get(db.Session, email)
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	if passwd == nil {
		fmt.Printf("error not found %v", email)
		return
	}
	passwd.Passwd = hashPass
	passwd.Enabled = true
	err = passwd.Insert(db.Session)
	fmt.Printf("err: %+v", err)
}

func createPasswd(email string, passwd string) {
	newpass, _ := hashPassword(passwd)
	newPasswd := models.UserPasswd{Email: email, Passwd: newpass, Enabled: true}

	err := newPasswd.Insert(db.Session)
	fmt.Printf("ads: %+v, err: %+v", newPasswd, err)

}
func (c *some) Exec() {
	// models.UserPasswd
	fmt.Printf("users Exec %+v", flag.Args())
	if flag.Args()[0] == "passwd" && len(flag.Args()) == 3 {
		updatePasswd(flag.Args()[1], flag.Args()[2])
		return
	}
	if flag.Args()[0] == "new" && len(flag.Args()) == 3 {
		createPasswd(flag.Args()[1], flag.Args()[2])
		return
	}
	c.Help()

}
