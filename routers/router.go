// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"gehpci/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebSocketController{},
			),
		),
		beego.NSNamespace("/proxys",
			beego.NSInclude(
				&controllers.ProxysController{},
			),
		),

		beego.NSNamespace("/command",
			beego.NSInclude(
				&controllers.CommandController{},
			),
		),
		beego.NSNamespace("/file",
			beego.NSInclude(
				&controllers.FileController{},
			),
		),
		beego.NSNamespace("/job",
			beego.NSInclude(
				&controllers.JobController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&controllers.AccountController{},
			),
		),
		beego.NSNamespace("/storage",
			beego.NSInclude(
				&controllers.StorageController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
