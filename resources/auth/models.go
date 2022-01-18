package auth

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordReq struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type AuthRes struct {
	Email      string            `json:"email"`
	Namespaces string            `json:"namespaces"`
	Token      string            `json:"token"`
	Attr       map[string]string `json:"attr"`
}
