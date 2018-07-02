package routers

import (
	"github.com/astaxie/beego"
	c "release/controllers"
)

func init() {
	beego.Router("/user/login", &c.Login{})
	beego.Router("/user/logout", &c.Logout{})
	beego.Router("/user/list", &c.UserList{})
	beego.Router("/user/add", &c.UserAdd{})
	beego.Router("/user/del/:username", &c.UserDel{})

	beego.Router("/release", &c.MainController{})
	beego.Router("/release/welcome", &c.Welcome{})
	beego.Router("/release/branch/list", &c.BranchList{})
	beego.Router("/release/branch/log", &c.BranchLog{})
	beego.Router("/release/branch/status", &c.BranchStatus{})
	beego.Router("/release/branch/pull/master", &c.BranchPull{})

	beego.Router("/release/order/list/:env/:page", &c.OrderList{})
	beego.Router("/release/order/list/false", &c.OrderListFalse{})
	beego.Router("/release/order/log/:timestamp", &c.OrderLog{})
	beego.Router("/release/order/search", &c.OrderSearch{})
	beego.Router("/release/order/add", &c.OrderAdd{})
	beego.Router("/release/order/modify/:id", &c.OrderModify{})
	beego.Router("/release/order/back/:timestamp", &c.OrderBack{})

	beego.Router("/release/pre/file", &c.PreFile{})
}
