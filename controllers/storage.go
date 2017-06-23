package controllers

import (
	"github.com/astaxie/beego"
)

// Operations about Storages (Not implemented)
type StorageController struct {
	beego.Controller
}

// @Title example
// @Description example.
// @Param       commad                body    string  true            "The commad json"
// @Success 200 {object} "models.example"
// @Failure 403 "example"
// @router / [post]
func (c *StorageController) About() {
	c.Data["json"] = "Not imp yet!"
	c.ServeJSON()
}
