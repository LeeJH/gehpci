package controllers

import (
	"encoding/json"
	"gehpci/models"
	"log"
	"net/http"
	"net/http/httputil"
)

type ProxysController struct {
	//	baseController
	BaseController
}

type myTransport struct {
	// Uncomment this if you want to capture the transport
	//CapturedTransport http.RoundTripper
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	// or, if you captured the transport
	//response, err := t.CapturedTransport.RoundTrip(request)

	// The httputil package provides a DumpResponse() func that will copy the
	// contents of the body into a []byte and return it. It also wraps it in an
	// ioutil.NopCloser and sets up the response to be passed on to the client.
	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		// copying the response body did not work
		return nil, err
	}

	// You may want to check the Content-Type header to decide how to deal with
	// the body. In this case, we're assuming it's text.
	log.Print(string(body))

	return response, err
}

// @Title build a proxy
// @Description
// @Param       machine            path    string  true            "The machine name"
// @Param	body	body	models.PortProxy	true	"The commad json"
// @Success 200 {object} models.PortProxy
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /:machine [post]
func (this *ProxysController) Proxys() {
	log.Printf("user: %#v", this.user)
	pp := models.NewPortProxy(this.user)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, pp)
	if err != nil {
		this.serveErrorCode(403, err) //serveError403(c, err)
		return
	}
	err = pp.Create()
	if err != nil {
		serveError403(&this.Controller, err)
	}
	this.Data["json"] = pp
	this.ServeJSON()
}

// @Title list live proxys
// @Description
// @Param       machine            path    string  true            "The machine name"
// @Success 200 {object} models.PortProxy
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /:machine [get]
func (this *ProxysController) ListProxys() {
	log.Printf("user: %#v", this.user)
	// pp := models.NewPortProxy(this.user)
	// err := json.Unmarshal(this.Ctx.Input.RequestBody, pp)
	// if err != nil {
	// 	this.serveErrorCode(403, err) //serveError403(c, err)
	// 	return
	// }
	// err = pp.Create()
	// if err != nil {
	// 	serveError403(&this.Controller, err)
	// }
	this.Data["json"] = models.ListPortProxys(this.user)
	this.ServeJSON()
}

// @Title Delete proxys
// @Description
// @Param       machine            path    string  true            "The machine name"
// @Param       proxyid            path    string  true            "The proxyid"
// @Success 200 {string} delete ok
// @Failure 403 {err} body is err info
// @Failure 401 need login
// @router /:machine/:proxyid [delete]
func (this *ProxysController) DeleteProxy() {
	log.Printf("user: %#v", this.user)
	proxyid := this.Ctx.Input.Param(":proxyid")
	log.Printf("try to del %s", proxyid)
	err := models.DeletePortProxy(proxyid, this.user)
	if err != nil {
		serveError403(&this.Controller, err)
	}
	this.Data["json"] = "delete port proxy ok "
	this.ServeJSON()
}
