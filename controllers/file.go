package controllers

import (
	"fmt"
	"gehpci/models"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Command
type FileController struct {
	BaseController //beego.Controller
}

// @Title Delete
// @Description Delete file or dir
// @Param       machine            path    string  true            "The machine name"
// @Param       pathname            path    string  true            "The path to list"
// @Param       all            query    bool  false            "if query this, will remove all"
// @Success 200 {string} result
// @Failure 403 {error} body is err info
// @Failure 401 need login
// @router /:machine/:pathname/* [delete]
func (c *FileController) Delete() {
	pathname, _, _ := getPath(&c.Controller)
	all := c.Ctx.Input.Query("all")
	var err error
	if all == "true" {
		err = os.RemoveAll(pathname)
	} else {
		err = os.Remove(pathname)
	}
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Error : " + err.Error()
		log.Printf("Error : %s\n", err.Error())
		c.ServeJSON()
		return
	}
	c.Data["json"] = "removed file/dir : " + pathname
	c.ServeJSON()
}

// ### /:machine/:pathname/*

// @Title List or Download(old version)
// @Description List file or dir stat
// @Param       machine            path    string  true            "The machine name"
// @Param       pathname            path    string  true            "The path to list"
// @Param       download            query    bool  false            "if query this, will download file"
// @Param       list            query    bool  false            "if query this, will list dir"
// @Param       simple            query    bool  false            "if query this,only return `name` and `isDir`"
// @Param       mkdir            query    bool  false            "if query this, will mkdir "
// @Success 200 {common.Result} result
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /:machine/?:pathname/* [get]
func (c *FileController) ListorDownload() {
	pathname, _, _ := getPath(&c.Controller)
	download := c.Ctx.Input.Query("download")
	list := c.Ctx.Input.Query("list")
	mkdir := c.Ctx.Input.Query("mkdir")
	var err error
	var fs models.FileStat
	if download == "true" {
		// do download
		err = downloadFile(&c.Controller, pathname)
		if err != nil {
			goto ListorDownloadError
		}
		return
	}
	if mkdir == "true" {
		err = os.MkdirAll(pathname, 0770)
		if err != nil {
			goto ListorDownloadError
		}
		c.Data["json"] = "create dir : " + pathname
		c.ServeJSON()
		return
	}
	// do list
	fs, err = models.GetFileStat(pathname)
	if err != nil {
		goto ListorDownloadError
	}
	if fs.Dir && list == "true" {
		fsl, err := models.GetDirList(pathname)
		if err != nil {
			goto ListorDownloadError
		}
		c.Data["json"] = fsl
		c.ServeJSON()
		return
	}
	c.Data["json"] = fs
	c.ServeJSON()
	return
ListorDownloadError:
	c.Ctx.Output.SetStatus(403)
	c.Data["json"] = "Error : " + err.Error()
	log.Printf("Error : %s\n", err.Error())
	c.ServeJSON()
	return
}

func getPath(c *beego.Controller) (pathname, machine string, err error) {
	//log.Printf("input : %#v \n", c.Ctx.Input)
	//log.Printf("machine : %#v \n", c.Ctx.Input.Param(":machine"))
	//log.Printf("download : %#v \n", c.Ctx.Input.Query("download"))
	filename := c.Ctx.Input.Param(":splat")
	machine = c.Ctx.Input.Param(":machine")
	pathname = c.Ctx.Input.Param(":pathname")
	pathname = path.Join("/", pathname, filename)
	if strings.HasSuffix(pathname, "/*") {
		pathname = pathname[0:(len(pathname) - 2)]
	}
	log.Println(pathname, machine, filename)
	return
}

func downloadFile(fc *beego.Controller, filename string) (err error) {
	log.Printf("try to download filename : %#v", filename)
	fc.Ctx.ResponseWriter.Header().Add("Content-Disposition", fmt.Sprintf("attachment; fileName=%s", path.Base(filename)))
	fc.Ctx.ResponseWriter.Header().Add("Content-Type", "multipart/form-data")
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	fs, err := os.Stat(filename)
	if err != nil {
		return
	}
	http.ServeContent(fc.Ctx.ResponseWriter, fc.Ctx.Request, path.Base(filename), fs.ModTime(), fd)
	return
}

/*// @Title List a dir
// @Description List file or dir stat
// @Param       machine            path    string  true            "The machine name"
// @Param       pathname            path    string  true            "The path to list"
// @Param       simple            query    bool  false            "if query this,only return `name` and `isDir`"
// @Success 200 {common.Result} result
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /list/:machine/:pathname/* [get]
func (c *FileController) ListDir() {
	pathname, _, _ := getPath(&c.Controller)
	fsl, err := models.GetDirList(pathname)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Error : " + err.Error()
		log.Printf("Error : %s\n", err.Error())
		c.ServeJSON()
		return
	}
	c.Data["json"] = fsl
	c.ServeJSON()
	return
}
*/

// @Title Upload
// @Description upload  file
// @Param       machine            path    string  true            "The machine name"
// @Param       pathname            path    string  false            "The path to list"
// @Param       body body bytes true            "The content of the file"
// @Success 200 {string} "upload ok"
// @Failure 403 "default errors"
// @Failure 401 "need login"
// @router /:machine/:pathname/* [put]
func (c *FileController) Upload() {
	filename, _, _ := getPath(&c.Controller)
	log.Print(c.Ctx.Input.Header("Content-Type"))
	if c.Ctx.Input.Header("Content-Type") == "application/x-www-form-urlencoded" {
		// #c.Ctx.Request.Body   // # = "application/octet-stream"
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Please set HTTP Header \"Content-Type\"=\"application/octet-stream\"  "
		c.ServeJSON()
		return
	}
	_, err := os.Stat(filename)
	if err == nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "file  exsit!"
		c.ServeJSON()
		return
	}
	err = nil
	_, err = os.Stat(path.Dir(filename))
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	dst, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	//log.Printf("%#v\n", c.Ctx.Request)
	wn, err := io.Copy(dst, c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = fmt.Sprintf("writen %d bytes.", wn)
	c.ServeJSON()
	return
}
