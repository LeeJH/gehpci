package models

import (
	"log"

	"github.com/astaxie/beego"
)

func init() {
	distmode := beego.AppConfig.String("distmode")
	log.Printf("distmode : %s", distmode)
	if distmode == "local" || distmode == "" {
		NewAuthMD = newAuthMD
	} else {
		authmode := beego.AppConfig.String(distmode + "::auth")
		log.Printf("authmode : %s", authmode)
		if authmode == "ldap" {
			NewAuthMD = newAuthMDLDAP
		} else {
			NewAuthMD = newAuthMD
		}

	}

	if NewAuthMD == nil {
		log.Fatal("Error: no available Auth model config.")
	}
}

type Author interface {
	authpw(string, string) bool
}

type AuthM struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Author
}

func NewAuthM(author Author) *AuthM {
	return &AuthM{Author: author}
}

func (authm *AuthM) Auth() bool {
	return authm.authpw(authm.Username, authm.Password)
}

var NewAuthMD func() *AuthM

// ************ authorm db

type AuthDataDeault struct {
}

func (au *AuthDataDeault) authpw(username, password string) bool {
	log.Printf("default login!")
	return Login(username, password)
}

func newAuthMD() *AuthM {
	return NewAuthM(&AuthDataDeault{})
}

// **********  authorm ldap
type AuthDataLDAP struct {
}

func (au *AuthDataLDAP) authpw(username, password string) bool {
	log.Printf("Error:Not implement yet! Auth LDAP")
	return false
}

func newAuthMDLDAP() *AuthM {
	return NewAuthM(&AuthDataLDAP{})
}
