package controllers

import (
	"github.com/astaxie/beego"
)

// Operations about Account,(Not implemented), such as look up Cputimes ,quotas ...
type AccountController struct {
	beego.Controller
}

// @Title account basic infos
// @Description search account basic infos.
// @Success 200 {object} "models.account.account"
// @Failure 403 "default error"
// @Failure 401 "login error"
// @router / [get]
func (c *AccountController) About() {
	c.Data["json"] = "Not imp yet!"
	c.ServeJSON()
}
