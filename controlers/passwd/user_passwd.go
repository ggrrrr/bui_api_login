package passwd

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/ggrrrr/bui_api_login/controlers"
	models "github.com/ggrrrr/bui_api_login/models/passwd"
	db "github.com/ggrrrr/bui_lib/db/cassandra"
)

// type AuthSystemErr error

func GetByEmail(c context.Context, email string) (*models.UserPasswd, *controlers.AuthError) {
	login, err := models.Get(db.Session, email)
	if err != nil {
		return nil, controlers.New(controlers.AUTH_ERR, err)
	}
	if login == nil {
		return nil, controlers.New(controlers.AUTH_NOT_FOUND, err)
	}
	return login, controlers.NewOK()
}

func CreateUserPasswd(c context.Context, o *models.UserPasswd) *controlers.AuthError {
	hash, err := hashPassword(o.Passwd)
	if err != nil {
		return controlers.Error(err)
	}
	o.Passwd = hash
	err = o.Insert(db.Session)
	return controlers.NewIfErr(err)
}

func ChangeUserPasswd(c context.Context, email, password, newPassword string) *controlers.AuthError {
	var err error
	login, err := models.Get(db.Session, email)
	if err != nil {
		return controlers.Error(err)
	}
	if login == nil {
		return controlers.New(controlers.AUTH_NOT_FOUND, nil)
	}
	if !checkPasswordHash(password, login.Passwd) {
		return controlers.New(controlers.AUTH_NOK, nil)
	}
	hash, err := hashPassword(newPassword)
	if err != nil {
		return controlers.Error(err)
	}
	newUserPassword := &models.UserPasswd{Email: email, Passwd: hash}
	// err = newUserPassword.Insert(db.Session)
	err = newUserPassword.UpdatePasswrd(db.Session)
	return controlers.NewIfErr(err)

}

func LockUserPasswd(c context.Context, email string) controlers.AuthError {
	return *controlers.ErrorStringf("not implemented")
}

func VerifyUserPasswd(c context.Context, email string, passwd string) (*models.UserPasswd, *controlers.AuthError) {
	var err error
	login, err := models.Get(db.Session, email)
	if err != nil {
		return nil, controlers.Error(err)
	}
	if login == nil {
		return nil, controlers.New(controlers.AUTH_NOT_FOUND, nil)
	}
	if !login.Enabled {
		return nil, controlers.New(controlers.AUTH_LOCKED, nil)
	}
	if checkPasswordHash(passwd, login.Passwd) {
		return login, controlers.NewOK()
	}
	return nil, controlers.New(controlers.AUTH_NOK, nil)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
