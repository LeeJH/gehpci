package controllers

import (
	"encoding/json"
	"fmt"
	"gehpci/models"
)

// Operations about Command
type JobController struct {
	BaseController //beego.Controller
}

// @Title Submit
// @Description Submit a HPCjob
// @Param	machine	path	string	true	"The machine name"
// @Param	body	body	models.HPCJob	true	"The job object"
// @Success 200 {string} result
// @Failure 403   err info
// @Failure 401 need login
// @router /:machine [post]
func (c *JobController) Submit() {
	job := &models.HPCJob{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, job)
	if err != nil {
		serveError403(&c.Controller, err)
		return
	}

	jobid, err := job.Submit()
	if err != nil {
		serveError403(&c.Controller, err)
		return
	}
	c.Data["json"] = struct {
		Jobid string `json:"jobid"`
	}{Jobid: fmt.Sprint(jobid)}
	c.ServeJSON()
}

// @Title Queue
// @Description Search Job in Queue
// @Param       machine            path    string  true            "The machine name"
// @Param       resource            query    bool  false            "if query this, will return resource info"
// @Success 200 {object} models.HPCJob
// @Failure 403 {error} body is err info
// @Failure 401 need login
// @router /:machine [get]
func (c *JobController) Queue() {
	resource := c.Ctx.Input.Query("resource")
	if resource == "true" {
		resl, err := models.ResourceInfo()
		if err != nil {
			//serveError403(&c.Controller, err)
			c.serveErrorCode(403, err)
			return
		}
		c.Data["json"] = resl
		c.ServeJSON()
		return
	}

	jobs, err := models.JobQueue()
	if err != nil {
		//serveError403(&c.Controller, err)
		c.serveErrorCode(403, err)
		return
	}
	c.Data["json"] = jobs
	c.ServeJSON()
}

// @Title JobInfo
// @Description Search Job Info
// @Param       machine            path    string  true            "The machine name"
// @Param       jobid            path    string  true            "The jobid"
// @Success 200 {object} models.HPCJob
// @Failure 403  err info
// @Failure 401 need login
// @router /:machine/:jobid [get]
func (c *JobController) JobInfo() {
	jobid := c.Ctx.Input.Param(":jobid")
	job, err := models.JobInfo(jobid)
	if err != nil {
		serveError403(&c.Controller, err)
		return
	}
	c.Data["json"] = job
	c.ServeJSON()
}

// @Title Delete
// @Description Delete a job
// @Param       machine            path    string  true            "The machine name"
// @Param       jobid            path    string  true            "The jobid"
// @Success 200 {string} result
// @Failure 403 {error} body is err info
// @Failure 401 need login
// @router /:machine/:jobid [delete]
func (c *JobController) Delete() {
	jobid := c.Ctx.Input.Param(":jobid")
	err := models.DeleteJob(jobid)
	if err != nil {
		serveError403(&c.Controller, err)
		return
	}
	c.Data["json"] = "Deleted job " + jobid
	c.ServeJSON()
}
