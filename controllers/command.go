package controllers

import (
	"encoding/json"
	"gehpci/models"
	"log"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Command
type CommandController struct {
	BaseController //beego.Controller
}
type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	sess := c.StartSession()
	username := sess.Get("username")
	if username == nil {
		c.Data["json"] = "not login"
		c.Ctx.Output.SetStatus(401)
		c.DestroySession()
		c.ServeJSON()
		c.StopRun()
	}
}

func serveError403(c *beego.Controller, err error) {
	c.Ctx.Output.SetStatus(403)
	c.Data["json"] = "Error : " + err.Error()
	log.Printf("Error : %s\n", err.Error())
	c.ServeJSON()
	return
}

func (c *BaseController) serveErrorCode(code int, err error) {
	c.Ctx.Output.SetStatus(code)
	c.Data["json"] = "Error : " + err.Error()
	log.Printf("Error : %s\n", err.Error())
	c.ServeJSON()
	c.StopRun()
	return
}

// @Title Run a Command
// @Description Run a command within limited time.
// @Param       machine            path    string  false            "The machine name"
// @Param	body	body	models.Command	true	"The commad json"
// @Success 200 {object} models.CommandResult
// @Failure 403 "Runtime error"
// @Failure 401 "need login"
// @router /run/?:machine [post]
func (c *CommandController) Run() {
	//realFunc := checkLogin(interRun)
	//realFunc(&c.Controller)
	//return
	//}
	//func interRun(c *beego.Controller) {
	var cmd = models.NewCommand()
	err := json.Unmarshal(c.Ctx.Input.RequestBody, cmd)
	if err != nil {
		c.serveErrorCode(403, err) //serveError403(c, err)
		return
	}
	result, err := cmd.Run()
	if err != nil {
		//serveError403(c, err)
		c.serveErrorCode(403, err)
		return
	}
	c.Data["json"] = result
	c.ServeJSON()
	return
}

// @Title Run a shell script
// @Description a simple command as 'bash -c "%command%" '
// @Param       command            formData    string  false            "The command script content,json format body also works"
// @Param       machine            path    string  false            "The machine name"
// @Success 200 {object} models.CommandResult
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /:machine [post]
func (c *CommandController) Shell() {
	//	realFunc := checkLogin(interShell)
	//	realFunc(&c.Controller)
	//	return
	//}
	//func interShell(c *beego.Controller) {

	mycmd := models.CommandShell{}
	if strings.HasPrefix(c.Ctx.Input.Header("Content-Type"), "application/json") {
		json.Unmarshal(c.Ctx.Input.RequestBody, &mycmd)
	} else {
		inputs := c.Input()
		mycmd.Command = inputs.Get("command")
	}

	if mycmd.Command == "" {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Error POST data format"
		c.ServeJSON()
		return
	}

	result, err := models.ShellRun(mycmd.Command)
	if err != nil {
		//serveError403(c, err)
		c.serveErrorCode(403, err)

		return
	}
	c.Data["json"] = result

	c.ServeJSON()
	return
}
