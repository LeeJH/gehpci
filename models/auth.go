package models

import (
	"log"

	"gehpci/extends/ldapc"

	"github.com/astaxie/beego"
)

func init() {
	distmode := beego.AppConfig.String("distmode")
	log.Printf("distmode : %s", distmode)
	if distmode == "local" || distmode == "" {
		//NewAuthMD = newAuthMD
		modelDAuthor = &AuthDataDeault{}
	} else {
		authmode := beego.AppConfig.String(distmode + "::auth")
		log.Printf("authmode : %s", authmode)
		if authmode == "ldap" {
			//NewAuthMD = newAuthMDLDAP
			modelDAuthor = newAuthLDAP()
		} else {
			//NewAuthMD = newAuthMD
			modelDAuthor = &AuthDataDeault{}

		}

	}

	if modelDAuthor == nil {
		log.Fatal("Error: no available Auth model config.")
	}
}

type Author interface {
	authpw(string, string) *User
}

var modelDAuthor Author

type AuthM struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Author
	user *User
}

func NewAuthM() *AuthM {
	return &AuthM{Author: modelDAuthor}
}

func (authm *AuthM) Auth() bool {
	var user *User
	if authm.Author != nil {
		user = authm.Author.authpw(authm.Username, authm.Password)
	} else {
		user = modelDAuthor.authpw(authm.Username, authm.Password)
	}
	if user == nil {
		return false
	}
	authm.user = user
	return true
}

func (authm *AuthM) GetUser() *User {
	return authm.user
}

//var NewAuthMD func() *AuthM

// ************ authorm db

type AuthDataDeault struct {
}

func (au *AuthDataDeault) authpw(username, password string) *User {
	log.Printf("default login!")
	user := LoginUser(username, password)
	return user
}

//func newAuthMD() *AuthM {
//	return NewAuthM(&AuthDataDeault{})
//}

// **********  authorm ldap
type AuthDataLDAP struct {
	*ldapc.Client
}

func (au *AuthDataLDAP) authpw(username, password string) *User {
	entry, _ := au.Authenticate(username, password)
	if entry != nil {
		log.Printf("Error:Not implement yet! Auth LDAP should pass")
	}
	return nil
}

func newAuthLDAP() *AuthDataLDAP {
	// read from config file
	host := beego.AppConfig.DefaultString("ldap::host", "localhost")
	port := beego.AppConfig.DefaultInt("ldap::port", 389)
	dn := beego.AppConfig.DefaultString("ldap::dn", "ou=people,dc=com")
	dn = "%s," + dn
	ldapclient := &ldapc.Client{
		Protocol:  ldapc.LDAP,
		Host:      host,
		Port:      port,
		TLSConfig: nil,
		Bind: &ldapc.DirectBind{
			UserDN: dn,
			Filter: "(&(objectClass=posixAccount)(uid=%s))",
		},
	}
	return &AuthDataLDAP{Client: ldapclient}
}

//func newAuthMDLDAP() *AuthM {
//	return NewAuthM(&AuthDataLDAP{})
//}
