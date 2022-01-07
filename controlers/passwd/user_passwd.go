package passwd

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/ggrrrr/bui_api_login/models"
	"github.com/ggrrrr/bui_lib/db"
)

// type AuthSystemErr error

type AUTH_RESULT byte

const AUTH_ERR AUTH_RESULT = 0
const AUTH_OK AUTH_RESULT = 1
const AUTH_NOK AUTH_RESULT = 2
const AUTH_NOT_FOUND AUTH_RESULT = 3
const AUTH_LOCKED AUTH_RESULT = 4

func AUTH_INFO(ok AUTH_RESULT) string {
	switch ok {
	case AUTH_ERR:
		return fmt.Sprintf("AUTH_ERR(%d)", ok)
	case AUTH_OK:
		return fmt.Sprintf("AUTH_OK(%d)", ok)
	case AUTH_NOK:
		return fmt.Sprintf("AUTH_NOK(%d)", ok)
	case AUTH_NOT_FOUND:
		return fmt.Sprintf("AUTH_NOT_FOUND(%d)", ok)
	case AUTH_LOCKED:
		return fmt.Sprintf("AUTH_LOCKED(%d)", ok)
	}
	return fmt.Sprintf("[%d]UNKOWN", ok)
}

func GetByEmail(c context.Context, email string) (AUTH_RESULT, *models.UserPasswd, error) {
	login, err := models.Get(db.Session, email)
	if err != nil {
		return AUTH_ERR, nil, err
	}
	if login == nil {
		return AUTH_NOT_FOUND, nil, nil
	}
	return AUTH_OK, login, nil
}

func CreateUserPasswd(c context.Context, o *models.UserPasswd) error {
	hash, err := hashPassword(o.Passwd)

	if err != nil {
		return err
	}
	o.Passwd = hash
	log.Printf("CreateLoginPasswd[%v]: %v", c, o)
	err = o.Insert(db.Session)
	return nil
}

func ChangeUserPasswd(c context.Context, o *models.UserPasswd) error {
	return nil
}

func LockUserPasswd(c context.Context, o *models.UserPasswd) error {
	return nil
}

func VerifyUserPasswd(c context.Context, email string, passwd string) (AUTH_RESULT, *models.UserPasswd, error) {
	var err error
	log.Printf("VerifyLoginPasswd[%+v]: user: %s, pass: %s", c, email, passwd)
	login, err := models.Get(db.Session, email)
	if err != nil {
		return AUTH_ERR, nil, err
	}
	if login == nil {
		return AUTH_NOT_FOUND, nil, nil
	}
	if checkPasswordHash(passwd, login.Passwd) {
		return AUTH_OK, login, nil
	}
	return AUTH_NOK, nil, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
