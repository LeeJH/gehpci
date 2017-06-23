package controllers

import (
	"encoding/json"
	"gehpci/models"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Auth
type AuthController struct {
	beego.Controller
}

// @Title Login
// @Description Logs into the system by username and password . got cookie if success.
// @Param       username                formData   string  true            "The username for login"
// @Param       password                formData    string  true            "The password for login"
// @Success 200 {string} "%username% login success"
// @Failure 403 "user not exist or password error"
// @router / [post]
func (c *AuthController) Login() {
	sess := c.StartSession()
	upo := models.NewAuthMD()
	if strings.HasPrefix(c.Ctx.Input.Header("Content-Type"), "application/json") {
		json.Unmarshal(c.Ctx.Input.RequestBody, upo)
	} else {
		inputs := c.Input()
		upo.Username, upo.Password = inputs.Get("username"), inputs.Get("password")
	}
	if upo.Auth() {
		sess.Set("username", upo.Username)
		c.Data["json"] = upo.Username + " login success"
	} else {
		c.Data["json"] = "user not exist"
		c.Ctx.Output.SetStatus(403)
		c.DestroySession()
	}
	c.ServeJSON()
}

// @Title Logout the system
// @Description the session willbe destory. cookie will be useless.
// @Success 200 {string} "logout success"
// @router / [delete]
func (c *AuthController) Logout() {
	//sess := c.StartSession()
	c.DestroySession()
	c.Data["json"] = "logout success"
	c.ServeJSON()
}

// @Title Status of auth
// @Description find if u are login status
// @Success 200 {string} username
// @Failure 403 "not login"
// @router / [get]
func (c *AuthController) Status() {
	sess := c.StartSession()
	username := sess.Get("username")
	if username == nil {
		c.Data["json"] = "not login"
		c.Ctx.Output.SetStatus(403)
		c.DestroySession()
	} else {
		c.Data["json"] = username
	}
	c.ServeJSON()
}
