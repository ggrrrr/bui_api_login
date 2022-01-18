package controlers

import "fmt"

type AUTH_RESULT byte

const AUTH_OK AUTH_RESULT = 0
const AUTH_ERR AUTH_RESULT = 1
const AUTH_NOK AUTH_RESULT = 2
const AUTH_NOT_FOUND AUTH_RESULT = 3
const AUTH_LOCKED AUTH_RESULT = 4

type AuthError struct {
	Result AUTH_RESULT
	Err    error
}

func New(ok AUTH_RESULT, err error) *AuthError {
	return &AuthError{Result: ok, Err: err}
}

func NewStringf(ok AUTH_RESULT, err string, a ...interface{}) *AuthError {
	return &AuthError{Result: ok, Err: fmt.Errorf(err, a...)}
}

func NewIfErr(err error) *AuthError {
	if err == nil {
		return &AuthError{Result: AUTH_OK, Err: nil}
	}
	return &AuthError{AUTH_ERR, err}
}

func NewOK() *AuthError {
	return &AuthError{AUTH_OK, nil}
}

func Unkown() *AuthError {
	return &AuthError{AUTH_ERR, fmt.Errorf("unkown")}
}

func Error(err error) *AuthError {
	return &AuthError{AUTH_ERR, err}
}

func ErrorStringf(err string, a ...interface{}) *AuthError {
	return &AuthError{AUTH_ERR, fmt.Errorf(err, a...)}
}

func (ok AuthError) String() string {
	switch ok.Result {
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
