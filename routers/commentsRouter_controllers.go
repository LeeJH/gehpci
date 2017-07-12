package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["gehpci/controllers:AccountController"] = append(beego.GlobalControllerRouter["gehpci/controllers:AccountController"],
		beego.ControllerComments{
			Method: "About",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:AuthController"] = append(beego.GlobalControllerRouter["gehpci/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:AuthController"] = append(beego.GlobalControllerRouter["gehpci/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:AuthController"] = append(beego.GlobalControllerRouter["gehpci/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Status",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:CommandController"] = append(beego.GlobalControllerRouter["gehpci/controllers:CommandController"],
		beego.ControllerComments{
			Method: "Run",
			Router: `/run/?:machine`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:CommandController"] = append(beego.GlobalControllerRouter["gehpci/controllers:CommandController"],
		beego.ControllerComments{
			Method: "Shell",
			Router: `/:machine`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:FileController"] = append(beego.GlobalControllerRouter["gehpci/controllers:FileController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:machine/:pathname/*`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:FileController"] = append(beego.GlobalControllerRouter["gehpci/controllers:FileController"],
		beego.ControllerComments{
			Method: "ListorDownload",
			Router: `/:machine/:pathname/*`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:FileController"] = append(beego.GlobalControllerRouter["gehpci/controllers:FileController"],
		beego.ControllerComments{
			Method: "Upload",
			Router: `/:machine/:pathname/*`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:JobController"] = append(beego.GlobalControllerRouter["gehpci/controllers:JobController"],
		beego.ControllerComments{
			Method: "Submit",
			Router: `/:machine`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:JobController"] = append(beego.GlobalControllerRouter["gehpci/controllers:JobController"],
		beego.ControllerComments{
			Method: "Queue",
			Router: `/:machine`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:JobController"] = append(beego.GlobalControllerRouter["gehpci/controllers:JobController"],
		beego.ControllerComments{
			Method: "JobInfo",
			Router: `/:machine/:jobid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:JobController"] = append(beego.GlobalControllerRouter["gehpci/controllers:JobController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:machine/:jobid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ObjectController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ObjectController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ObjectController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ObjectController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ObjectController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ProxysController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ProxysController"],
		beego.ControllerComments{
			Method: "Proxys",
			Router: `/:machine`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ProxysController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ProxysController"],
		beego.ControllerComments{
			Method: "ListProxys",
			Router: `/:machine`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:ProxysController"] = append(beego.GlobalControllerRouter["gehpci/controllers:ProxysController"],
		beego.ControllerComments{
			Method: "DeleteProxy",
			Router: `/:machine/:proxyid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:StorageController"] = append(beego.GlobalControllerRouter["gehpci/controllers:StorageController"],
		beego.ControllerComments{
			Method: "About",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:UserController"] = append(beego.GlobalControllerRouter["gehpci/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["gehpci/controllers:WebSocketController"] = append(beego.GlobalControllerRouter["gehpci/controllers:WebSocketController"],
		beego.ControllerComments{
			Method: "BindBash",
			Router: `/bindbash/:machine`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
